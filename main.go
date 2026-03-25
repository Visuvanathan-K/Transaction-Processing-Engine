package main

import (
	"log"
	"net/http"

	"transaction-engine/handlers"
	"transaction-engine/store"
)

func main() {
	// Initialize in-memory store with seed data
	store.Init()

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/transaction", handlers.HandleTransaction)
	mux.HandleFunc("/api/card/balance/", handlers.HandleGetBalance)
	mux.HandleFunc("/api/card/transactions/", handlers.HandleGetTransactions)

	log.Println("Transaction Processing Engine running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
