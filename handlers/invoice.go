package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type CreateInvoiceRequest struct {
	MerchantID  string `json:"merchant_id"`
	CustomerID  string `json:"customer_id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

func isUpperAlpha(s string) bool {
	if len(s) != 3 {
		return false
	}
	for _, c := range s {
		if !unicode.IsUpper(c) || !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

func CreateInvoiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("Received POST /invoices", ctx)

		var req CreateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Basic input validation
		if req.Amount <= 0 {
			http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
			return
		}
		if !isUpperAlpha(req.Currency) {
			http.Error(w, "Currency must be a 3-letter uppercase code (e.g. USD)", http.StatusBadRequest)
			return
		}

		invoiceID := "inv_" + uuid.New().String()[:8]
		createdAt := time.Now().Format(time.RFC3339)

		_, err := db.Exec(
			`INSERT INTO invoices (id, merchant_id, customer_id, amount, currency, description, status, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, 'unpaid', ?)`,
			invoiceID, req.MerchantID, req.CustomerID, req.Amount, req.Currency, req.Description, createdAt,
		)
		if err != nil {
			log.Printf("Insert error: %v\n", err)
			http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"invoice_id": invoiceID,
			"created_at": createdAt,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func GetInvoiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("Received GET /invoices/{id}", ctx)

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 || parts[2] == "" {
			http.Error(w, "Missing invoice ID", http.StatusBadRequest)
			return
		}
		invoiceID := parts[2]

		row := db.QueryRow(`SELECT id, merchant_id, customer_id, amount, currency, description, status, created_at FROM invoices WHERE id = ?`, invoiceID)

		var inv struct {
			ID          string `json:"id"`
			MerchantID  string `json:"merchant_id"`
			CustomerID  string `json:"customer_id"`
			Amount      int    `json:"amount"`
			Currency    string `json:"currency"`
			Description string `json:"description"`
			Status      string `json:"status"`
			CreatedAt   string `json:"created_at"`
		}

		err := row.Scan(&inv.ID, &inv.MerchantID, &inv.CustomerID, &inv.Amount, &inv.Currency, &inv.Description, &inv.Status, &inv.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Invoice not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch invoice", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inv)
	}
}
