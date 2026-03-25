package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"transaction-engine/models"
	"transaction-engine/store"
)

// generateID creates a random hex transaction ID without external dependencies
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// --- POST /api/transaction ---

func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.TransactionResponse{
			Status:   "FAILED",
			RespCode: "99",
			Message:  "Invalid request body",
		})
		return
	}

	// --- Validate card exists ---
	card, exists := store.GetCard(req.CardNumber)
	if !exists {
		logTxFull(req.CardNumber, models.TxWithdraw, req.Amount, models.TxFailed)
		writeJSON(w, http.StatusOK, models.TransactionResponse{
			Status:   "FAILED",
			RespCode: "05",
			Message:  "Invalid card",
		})
		return
	}

	// --- Validate card is ACTIVE ---
	if card.Status != models.StatusActive {
		logTxFull(req.CardNumber, models.TxWithdraw, req.Amount, models.TxFailed)
		writeJSON(w, http.StatusOK, models.TransactionResponse{
			Status:   "FAILED",
			RespCode: "05",
			Message:  "Card is blocked",
		})
		return
	}

	// --- Validate PIN (compare SHA-256 hashes) ---
	if store.HashPIN(req.PIN) != card.PinHash {
		logTxFull(req.CardNumber, models.TxWithdraw, req.Amount, models.TxFailed)
		writeJSON(w, http.StatusOK, models.TransactionResponse{
			Status:   "FAILED",
			RespCode: "06",
			Message:  "Invalid PIN",
		})
		return
	}

	// --- Validate transaction type ---
	txType := models.TransactionType(strings.ToLower(req.Type))
	if txType != models.TxWithdraw && txType != models.TxTopup {
		writeJSON(w, http.StatusOK, models.TransactionResponse{
			Status:   "FAILED",
			RespCode: "99",
			Message:  "Invalid transaction type. Allowed: withdraw, topup",
		})
		return
	}

	// --- Validate amount ---
	if req.Amount <= 0 {
		writeJSON(w, http.StatusOK, models.TransactionResponse{
			Status:   "FAILED",
			RespCode: "99",
			Message:  "Amount must be greater than 0",
		})
		return
	}

	// --- Process transaction ---
	var delta float64
	if txType == models.TxWithdraw {
		if card.Balance < req.Amount {
			logTxFull(req.CardNumber, txType, req.Amount, models.TxFailed)
			writeJSON(w, http.StatusOK, models.TransactionResponse{
				Status:   "FAILED",
				RespCode: "99",
				Message:  "Insufficient balance",
			})
			return
		}
		delta = -req.Amount
	} else {
		delta = req.Amount
	}

	newBalance := store.UpdateBalance(req.CardNumber, delta)
	logTxFull(req.CardNumber, txType, req.Amount, models.TxSuccess)

	writeJSON(w, http.StatusOK, models.TransactionResponse{
		Status:   "SUCCESS",
		RespCode: "00",
		Balance:  newBalance,
	})
}

// --- GET /api/card/balance/{cardNumber} ---

func HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := strings.TrimPrefix(r.URL.Path, "/api/card/balance/")
	if cardNumber == "" {
		http.Error(w, "Card number required", http.StatusBadRequest)
		return
	}

	card, exists := store.GetCard(cardNumber)
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Card not found"})
		return
	}

	writeJSON(w, http.StatusOK, models.BalanceResponse{
		CardNumber: card.CardNumber,
		CardHolder: card.CardHolder,
		Balance:    card.Balance,
		Status:     string(card.Status),
	})
}

// --- GET /api/card/transactions/{cardNumber} ---

func HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := strings.TrimPrefix(r.URL.Path, "/api/card/transactions/")
	if cardNumber == "" {
		http.Error(w, "Card number required", http.StatusBadRequest)
		return
	}

	if _, exists := store.GetCard(cardNumber); !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Card not found"})
		return
	}

	txs := store.GetTransactions(cardNumber)
	writeJSON(w, http.StatusOK, txs)
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func logTxFull(cardNumber string, txType models.TransactionType, amount float64, status models.TxStatus) {
	store.AddTransaction(models.Transaction{
		TransactionID: fmt.Sprintf("%s", generateID()),
		CardNumber:    cardNumber,
		Type:          txType,
		Amount:        amount,
		Status:        status,
		Timestamp:     time.Now(),
	})
}
