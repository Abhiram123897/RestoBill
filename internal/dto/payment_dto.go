package dto

// RecordPaymentRequest is the request body for recording a payment
type RecordPaymentRequest struct {
	BillID    string  `json:"bill_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Method    string  `json:"method" binding:"required"`
	Reference string  `json:"reference"`
}

// PaymentResponse is the API response for a payment
type PaymentResponse struct {
	ID        string  `json:"id"`
	BillID    string  `json:"bill_id"`
	Amount    float64 `json:"amount"`
	Method    string  `json:"method"`
	Reference string  `json:"reference"`
	CreatedAt string  `json:"created_at"`
}
