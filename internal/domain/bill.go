package domain

import "time"

// Bill represents a generated bill for an order
type Bill struct {
	ID             string    `json:"id"`
	OrderID        string    `json:"order_id"`
	Subtotal       float64   `json:"subtotal"`
	TaxAmount      float64   `json:"tax_amount"`
	TaxRate        float64   `json:"tax_rate"`
	DiscountAmount float64   `json:"discount_amount"`
	DiscountReason string    `json:"discount_reason"`
	Total          float64   `json:"total"`
	Status         string    `json:"status"`
	GeneratedBy    string    `json:"generated_by"`
	CreatedAt      time.Time `json:"created_at"`
}
