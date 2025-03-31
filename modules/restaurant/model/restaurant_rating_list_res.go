package restaurantmodel

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantRatingListRes struct {
	RestaurantId   uuid.UUID  `json:"restaurantId"`
	RestaurantName string     `json:"restaurantName"`
	UserId         string     `json:"userId"`
	UserName       string     `json:"userName"`
	Avatar         *string    `json:"avatar"`
	Rating         float64    `json:"rating"`
	Comment        *string    `json:"comment"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}
