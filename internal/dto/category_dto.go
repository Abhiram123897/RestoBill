package dto

// CreateCategoryRequest is the request body for creating a category
type CreateCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	DisplayOrder int    `json:"display_order"`
}

// UpdateCategoryRequest is the request body for updating a category
type UpdateCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	DisplayOrder int    `json:"display_order"`
}

// CategoryResponse is the API response for a category
type CategoryResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
