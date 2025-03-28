package restaurantmodel

import (
	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantRating struct {
	ID           uuid.UUID `gorm:"column:id"`
	UserID       string    `gorm:"column:user_id"`
	RestaurantID string    `gorm:"column:restaurant_id"`
	Point        float64   `gorm:"column:point"`
	Comment      *string   `gorm:"column:comment"`
	Status       string    `gorm:"column:status"`
	sharedmodel.AbstractInfo
}
