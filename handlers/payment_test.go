package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Helper: Insert a test invoice into the test DB
func insertTestInvoice(db *sql.DB, id, status string) {
	db.Exec(`INSERT INTO invoices (id, merchant_id, customer_id, amount, currency, description, status, created_at)
	VALUES (?, 'merchant', 'customer', 10, 'USD', 'desc', ?, ?)`,
		id, status, time.Now().Format(time.RFC3339))
}

func setupPaymentTestDB() *sql.DB {
	db, _ := sql.Open("sqlite3", ":memory:")
	db.Exec(`
	CREATE TABLE IF NOT EXISTS invoices (
		id TEXT PRIMARY KEY,
		merchant_id TEXT NOT NULL,
		customer_id TEXT NOT NULL,
		amount INTEGER NOT NULL,
		currency TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'unpaid',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	db.Exec(`
	CREATE TABLE IF NOT EXISTS payments (
		id TEXT PRIMARY KEY,
		invoice_id TEXT NOT NULL,
		method TEXT NOT NULL,
		status TEXT NOT NULL,
		paid_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (invoice_id) REFERENCES invoices(id)
	);`)
	return db
}

func TestProcessPaymentHandler(t *testing.T) {
	db := setupPaymentTestDB()
	handler := ProcessPaymentHandler(db)

	// Try paying a non-existent invoice (should 404)
	body := []byte(`{"invoice_id":"no_such_id", "method":"card", "details":{"card_number":"4111111111111111","expiry":"12/25","cvv":"123"}}`)
	req := httptest.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent invoice, got %d", w.Code)
	}

	// Insert a paid invoice and try paying it (should 400)
	insertTestInvoice(db, "paid_inv", "paid")
	body = []byte(`{"invoice_id":"paid_inv", "method":"card", "details":{"card_number":"4111111111111111","expiry":"12/25","cvv":"123"}}`)
	req = httptest.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for already paid invoice, got %d", w.Code)
	}

	// Insert a valid unpaid invoice, try paying (randomly may fail, so try several times)
	insertTestInvoice(db, "ok_inv", "unpaid")
	success := false
	for i := 0; i < 10; i++ {
		body = []byte(`{"invoice_id":"ok_inv", "method":"card", "details":{"card_number":"4111111111111111","expiry":"12/25","cvv":"123"}}`)
		req = httptest.NewRequest("POST", "/payments", bytes.NewBuffer(body))
		w = httptest.NewRecorder()
		handler(w, req)
		if w.Code == http.StatusOK {
			success = true
			var resp map[string]string
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["status"] != "paid" {
				t.Errorf("expected status 'paid', got %v", resp["status"])
			}
			break
		}
	}
	if !success {
		t.Errorf("payment never succeeded after 10 tries (random failure sim)")
	}
}
