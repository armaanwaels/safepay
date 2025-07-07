package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() *sql.DB {
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
	return db
}

func TestCreateInvoiceHandler_Validation(t *testing.T) {
	db := setupTestDB()
	handler := CreateInvoiceHandler(db)

	// Try amount = 0 (should fail)
	body := []byte(`{"merchant_id":"a","customer_id":"b","amount":0,"currency":"USD","description":"x"}`)
	req := httptest.NewRequest("POST", "/invoices", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero amount, got %d", w.Code)
	}

	// Try invalid currency (should fail)
	body = []byte(`{"merchant_id":"a","customer_id":"b","amount":1,"currency":"usd","description":"x"}`)
	req = httptest.NewRequest("POST", "/invoices", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad currency, got %d", w.Code)
	}

	// Try valid invoice (should succeed)
	body = []byte(`{"merchant_id":"a","customer_id":"b","amount":1,"currency":"USD","description":"x"}`)
	req = httptest.NewRequest("POST", "/invoices", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for valid invoice, got %d", w.Code)
	}
}
