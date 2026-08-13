package domain

import "time"

// Category represents a menu category in the restaurant
type Category struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
