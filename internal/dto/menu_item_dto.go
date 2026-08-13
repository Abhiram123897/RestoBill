package dto

// CreateMenuItemRequest is the request body for creating a menu item
type CreateMenuItemRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	IsAvailable *bool   `json:"is_available"`
}

// UpdateMenuItemRequest is the request body for updating a menu item
type UpdateMenuItemRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	IsAvailable *bool   `json:"is_available"`
}

// MenuItemResponse is the API response for a menu item
type MenuItemResponse struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
