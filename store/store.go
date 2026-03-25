package store

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"transaction-engine/models"
)

var (
	mu           sync.RWMutex
	cards        = make(map[string]*models.Card)
	transactions = make(map[string][]models.Transaction) // cardNumber -> []Transaction
)

// hashPIN returns the SHA-256 hex string of a plain PIN
func HashPIN(pin string) string {
	h := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", h)
}

// Init seeds the in-memory store with example cards
func Init() {
	seed := []struct {
		number string
		holder string
		pin    string
		balance float64
		status models.CardStatus
	}{
		{"4123456789012345", "John Doe", "1234", 1000.00, models.StatusActive},
		{"4123456789012346", "Jane Smith", "5678", 500.00, models.StatusActive},
		{"4123456789012347", "Bob Johnson", "9999", 250.00, models.StatusBlocked},
	}

	for _, s := range seed {
		cards[s.number] = &models.Card{
			CardNumber: s.number,
			CardHolder: s.holder,
			PinHash:    HashPIN(s.pin),
			Balance:    s.balance,
			Status:     s.status,
		}
	}
}

// GetCard retrieves a card by number (read lock)
func GetCard(cardNumber string) (*models.Card, bool) {
	mu.RLock()
	defer mu.RUnlock()
	c, ok := cards[cardNumber]
	return c, ok
}

// UpdateBalance updates a card's balance (write lock) and returns the new balance
func UpdateBalance(cardNumber string, delta float64) float64 {
	mu.Lock()
	defer mu.Unlock()
	cards[cardNumber].Balance += delta
	return cards[cardNumber].Balance
}

// AddTransaction appends a transaction log entry
func AddTransaction(tx models.Transaction) {
	mu.Lock()
	defer mu.Unlock()
	transactions[tx.CardNumber] = append(transactions[tx.CardNumber], tx)
}

// GetTransactions returns all transactions for a card
func GetTransactions(cardNumber string) []models.Transaction {
	mu.RLock()
	defer mu.RUnlock()
	txs := transactions[cardNumber]
	if txs == nil {
		return []models.Transaction{}
	}
	return txs
}
