package restaurantmodel

import "github.com/google/uuid"

type RestaurantInfoDto struct {
	Restaurant
	AvgPoint   float64    `json:"avgPoint"`
	CommentQty int        `json:"commentQty"`
	LikesQty   int        `json:"likesQty"`
	FoodInfos  []FoodInfo `json:"foodInfos"`
}

type FoodInfo struct {
	FoodId     uuid.UUID `json:"food_id"`
	FoodName   string    `json:"food_name"`
	CategoryId uuid.UUID `json:"category_id"`
}
