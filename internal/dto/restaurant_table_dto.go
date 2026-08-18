package dto

// CreateRestaurantTableRequest is the request body for creating a table
type CreateRestaurantTableRequest struct {
	Number   int `json:"number" binding:"required"`
	Capacity int `json:"capacity" binding:"required"`
}

// UpdateRestaurantTableRequest is the request body for updating a table
type UpdateRestaurantTableRequest struct {
	Number   int `json:"number" binding:"required"`
	Capacity int `json:"capacity" binding:"required"`
}

// UpdateTableStatusRequest is the request body for updating table status
type UpdateTableStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// RestaurantTableResponse is the API response for a table
type RestaurantTableResponse struct {
	ID        string `json:"id"`
	Number    int    `json:"number"`
	Capacity  int    `json:"capacity"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
