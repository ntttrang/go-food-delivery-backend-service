package foodmodel

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	Id           uuid.UUID  `json:"id"`
	RestaurantId uuid.UUID  `json:"restaurant_id"`
	CategoryId   uuid.UUID  `json:"category_id,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	Price        float64    `json:"price"`
	Images       string     `json:"images"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (Food) TableName() string {
	return "foods"
}
