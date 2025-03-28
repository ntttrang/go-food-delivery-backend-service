package restaurantmodel

import "github.com/google/uuid"

type RestaurantFoodDto struct {
	FoodID uuid.UUID `json:"foodId"`
	Status string    `json:"status"`
}
