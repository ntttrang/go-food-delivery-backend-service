package restaurantmodel

import (
	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantLike struct {
	RestaurantID uuid.UUID `gorm:"column:restaurant_id"`
	UserID       uuid.UUID `gorm:"column:user_id"`
	sharedmodel.AbstractInfo
}
