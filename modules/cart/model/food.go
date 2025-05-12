package cartmodel

import "github.com/google/uuid"

type Food struct {
	Id           uuid.UUID `gorm:"column:food_id" json:"id"`
	RestaurantId uuid.UUID `gorm:"column:food_id" json:"restaurantId"`
	CategoryId   uuid.UUID `json:"categoryId"`
	Name         string
	Description  string
	Images       string
	Price        float64
}
