package foodmodel

import (
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ListFoodRes struct {
	Items      []FoodSearchResDto    `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}
