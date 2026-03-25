package models

import "time"

type CardStatus string

const (
	StatusActive  CardStatus = "ACTIVE"
	StatusBlocked CardStatus = "BLOCKED"
)

type TransactionType string

const (
	TxWithdraw TransactionType = "withdraw"
	TxTopup    TransactionType = "topup"
)

type TxStatus string

const (
	TxSuccess TxStatus = "SUCCESS"
	TxFailed  TxStatus = "FAILED"
)

// Card represents a stored card record
type Card struct {
	CardNumber string
	CardHolder string
	PinHash    string // SHA-256 hex of PIN
	Balance    float64
	Status     CardStatus
}

// Transaction represents a transaction log record
type Transaction struct {
	TransactionID string
	CardNumber    string
	Type          TransactionType
	Amount        float64
	Status        TxStatus
	Timestamp     time.Time
}

// TransactionRequest is the incoming request body
type TransactionRequest struct {
	CardNumber string  `json:"cardNumber"`
	PIN        string  `json:"pin"`
	Type       string  `json:"type"`
	Amount     float64 `json:"amount"`
}

// TransactionResponse is returned after processing
type TransactionResponse struct {
	Status  string  `json:"status"`
	RespCode string `json:"respCode"`
	Balance float64 `json:"balance,omitempty"`
	Message string  `json:"message,omitempty"`
}

// BalanceResponse is returned for GET balance
type BalanceResponse struct {
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	Balance    float64 `json:"balance"`
	Status     string  `json:"status"`
}
