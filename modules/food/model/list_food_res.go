package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ListFoodRes struct {
	Items      []FoodSearchResDto    `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

type FoodSearchResDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	sharedModel.DateDto
}
