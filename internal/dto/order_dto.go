package dto

// CreateOrderRequest is the request body for creating an order
type CreateOrderRequest struct {
	TableID   *string `json:"table_id"`
	OrderType string  `json:"order_type" binding:"required"`
	CreatedBy string  `json:"created_by" binding:"required"`
}

// AddOrderItemRequest is the request body for adding an item to an order
type AddOrderItemRequest struct {
	MenuItemID string `json:"menu_item_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	Notes      string `json:"notes"`
}

// UpdateOrderStatusRequest is the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// OrderItemResponse is the API response for an order item
type OrderItemResponse struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	MenuItemID string  `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	Notes      string  `json:"notes"`
	CreatedAt  string  `json:"created_at"`
}

// OrderResponse is the API response for an order
type OrderResponse struct {
	ID        string              `json:"id"`
	TableID   *string             `json:"table_id"`
	Status    string              `json:"status"`
	OrderType string              `json:"order_type"`
	CreatedBy string              `json:"created_by"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}
