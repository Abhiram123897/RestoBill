package domain

import "time"

// Payment represents a payment recorded against a bill
type Payment struct {
	ID        string    `json:"id"`
	BillID    string    `json:"bill_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
}
