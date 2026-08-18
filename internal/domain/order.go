package domain

import "time"

// Order represents a customer order in the restaurant
type Order struct {
	ID        string    `json:"id"`
	TableID   *string   `json:"table_id"`
	Status    string    `json:"status"`
	OrderType string    `json:"order_type"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ID         string    `json:"id"`
	OrderID    string    `json:"order_id"`
	MenuItemID string    `json:"menu_item_id"`
	Quantity   int       `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}
