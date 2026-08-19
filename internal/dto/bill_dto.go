package dto

// GenerateBillRequest is the request body for generating a bill
type GenerateBillRequest struct {
	OrderID        string  `json:"order_id" binding:"required"`
	TaxRate        float64 `json:"tax_rate" binding:"required"`
	DiscountAmount float64 `json:"discount_amount"`
	DiscountReason string  `json:"discount_reason"`
	GeneratedBy    string  `json:"generated_by" binding:"required"`
}

// ApplyDiscountRequest is the request body for applying discount to a bill
type ApplyDiscountRequest struct {
	DiscountAmount float64 `json:"discount_amount" binding:"required"`
	DiscountReason string  `json:"discount_reason" binding:"required"`
}

// BillResponse is the API response for a bill
type BillResponse struct {
	ID             string  `json:"id"`
	OrderID        string  `json:"order_id"`
	Subtotal       float64 `json:"subtotal"`
	TaxAmount      float64 `json:"tax_amount"`
	TaxRate        float64 `json:"tax_rate"`
	DiscountAmount float64 `json:"discount_amount"`
	DiscountReason string  `json:"discount_reason"`
	Total          float64 `json:"total"`
	Status         string  `json:"status"`
	GeneratedBy    string  `json:"generated_by"`
	CreatedAt      string  `json:"created_at"`
}
