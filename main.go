package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"safepay/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// extra logging setup (file or console)
	log.SetOutput(os.Stdout)
	log.Println("Safepay server starting...")

	// Open SQLite DB (creates file if not exists)
	db, err := sql.Open("sqlite3", "./safepay.db")
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	// Run schema migration
	schema, err := ioutil.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema.sql:", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to execute schema.sql:", err)
	}
	fmt.Println("Migration ran successfully.")

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/invoices", handlers.CreateInvoiceHandler(db))     // POST /invoices
	mux.HandleFunc("/invoices/", handlers.GetInvoiceHandler(db))       // GET /invoices/{id}
	mux.HandleFunc("/payments", handlers.ProcessPaymentHandler(db))    // POST /payments
	mux.HandleFunc("/payments/", handlers.GetPaymentStatusHandler(db)) // GET /payments/{invoice_id}

	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
