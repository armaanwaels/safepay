package models

type Invoice struct {
	ID          string `json:"id"`
	MerchantID  string `json:"merchant_id"`
	CustomerID  string `json:"customer_id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}
