package restaurantmodel

import (
	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantFood struct {
	RestaurantId uuid.UUID `gorm:"column:restaurant_id"`
	FoodId       uuid.UUID `gorm:"column:food_id"`
	Status       string    `gorm:"column:status"`
	sharedmodel.DateDto
}

func (RestaurantFood) TableName() string {
	return "restaurant_foods"
}
