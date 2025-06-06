package restaurantmodel

import "github.com/google/uuid"

type Foods struct {
	Id           uuid.UUID
	Name         string
	Description  string
	Images       string
	Price        float64
	AvgPoint     float64
	CommentQty   int
	CategoryId   uuid.UUID
	RestaurantId uuid.UUID
	Status       string
}
