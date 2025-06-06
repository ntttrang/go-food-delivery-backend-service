package cartmodel

import "github.com/google/uuid"

type Food struct {
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
