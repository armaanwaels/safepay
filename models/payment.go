package models

type Payment struct {
	ID        string `json:"payment_id"`
	InvoiceID string `json:"invoice_id"`
	Method    string `json:"method"`
	Status    string `json:"status"`
	PaidAt    string `json:"paid_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
