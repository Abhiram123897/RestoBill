package domain

import "time"

// MenuItem represents a dish/drink on the restaurant menu
type MenuItem struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"category_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
