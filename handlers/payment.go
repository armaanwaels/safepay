package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type ProcessPaymentRequest struct {
	InvoiceID string `json:"invoice_id"`
	Method    string `json:"method"`
	Details   struct {
		CardNumber string `json:"card_number"`
		Expiry     string `json:"expiry"`
		CVV        string `json:"cvv"`
	} `json:"details"`
}

func ProcessPaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("Received POST /payments", ctx)

		var req ProcessPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Input validation for card payments
		if req.Method == "card" {
			num := strings.ReplaceAll(req.Details.CardNumber, " ", "")
			if len(num) != 16 {
				http.Error(w, "Card number must be 16 digits", http.StatusBadRequest)
				return
			}
			// Simple expiry check (mm/yy or mm/yyyy)
			if !(len(req.Details.Expiry) == 5 || len(req.Details.Expiry) == 7) {
				http.Error(w, "Expiry must be in MM/YY or MM/YYYY format", http.StatusBadRequest)
				return
			}
			if len(req.Details.CVV) < 3 {
				http.Error(w, "CVV must be at least 3 digits", http.StatusBadRequest)
				return
			}
		}

		// Check if invoice exists and is unpaid
		var status string
		err := db.QueryRow("SELECT status FROM invoices WHERE id = ?", req.InvoiceID).Scan(&status)
		if err == sql.ErrNoRows {
			http.Error(w, "Invoice not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		if status == "paid" {
			http.Error(w, "Invoice already paid", http.StatusBadRequest)
			return
		}

		// Randomly simulate payment failure (20% chance)
		if rand.Float32() < 0.2 {
			http.Error(w, "Payment failed (random simulation)", http.StatusPaymentRequired)
			return
		}

		paymentID := "pay_" + RandString8()
		paidAt := time.Now().Format(time.RFC3339)
		payStatus := "paid"

		_, err = db.Exec(
			`INSERT INTO payments (id, invoice_id, method, status, paid_at)
             VALUES (?, ?, ?, ?, ?)`,
			paymentID, req.InvoiceID, req.Method, payStatus, paidAt,
		)
		if err != nil {
			log.Printf("Payment insert error: %v\n", err)
			http.Error(w, "Failed to process payment", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`UPDATE invoices SET status = 'paid' WHERE id = ?`, req.InvoiceID)
		if err != nil {
			log.Printf("Invoice update error: %v\n", err)
			http.Error(w, "Failed to update invoice status", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"payment_id": paymentID,
			"status":     payStatus,
			"paid_at":    paidAt,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// Helper for random payment IDs
func RandString8() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetPaymentStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("Received GET /payments/{invoice_id}", ctx)

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 || parts[2] == "" {
			http.Error(w, "Missing invoice ID", http.StatusBadRequest)
			return
		}
		invoiceID := parts[2]

		row := db.QueryRow(
			`SELECT id, status, paid_at FROM payments WHERE invoice_id = ? ORDER BY paid_at DESC LIMIT 1`,
			invoiceID,
		)

		var paymentID, status, paidAt string
		err := row.Scan(&paymentID, &status, &paidAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Payment not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"invoice_id": invoiceID,
			"payment_id": paymentID,
			"status":     status,
			"paid_at":    paidAt,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
