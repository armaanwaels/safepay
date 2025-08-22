package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"safepay/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Setup logging
	log.SetOutput(os.Stdout)
	log.Println("Safepay server starting...")

	// Get DB path from environment (default: ./safepay.db)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./safepay.db"
	}
	log.Printf("Using database file: %s\n", dbPath)

	// Open SQLite DB (creates file if not exists)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	// Run schema migration
	schemaPath := "./db/schema.sql"
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", schemaPath, err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("Failed to execute schema migration: %v", err)
	}
	fmt.Println("Migration ran successfully.")

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/invoices", handlers.CreateInvoiceHandler(db))     // POST /invoices
	mux.HandleFunc("/invoices/", handlers.GetInvoiceHandler(db))       // GET /invoices/{id}
	mux.HandleFunc("/payments", handlers.ProcessPaymentHandler(db))    // POST /payments
	mux.HandleFunc("/payments/", handlers.GetPaymentStatusHandler(db)) // GET /payments/{invoice_id}

	// Start server
	addr := ":8080"
	log.Printf("Server running at http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
