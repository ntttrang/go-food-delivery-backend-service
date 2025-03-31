package restaurantmodel

import (
	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantFood struct {
	RestaurantID uuid.UUID `gorm:"column:restaurant_id"`
	FoodID       uuid.UUID `gorm:"column:food_id"`
	Status       string    `gorm:"column:status"`
	sharedmodel.AbstractInfo
}

func (r RestaurantFoodDto) ConvertToRestaurantFood() *RestaurantFood {
	return &RestaurantFood{
		FoodID: r.FoodID,
		Status: r.Status,
	}
}
