package domain

import "time"

// RestaurantTable represents a physical table in the restaurant
type RestaurantTable struct {
	ID        string    `json:"id"`
	Number    int       `json:"number"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
