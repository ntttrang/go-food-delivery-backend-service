package foodmodel

import (
	"time"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FoodRatingListRes struct {
	Items      []FoodRatingListDto   `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

type FoodRatingListDto struct {
	Id        uuid.UUID  `json:"id"`
	FoodId    uuid.UUID  `json:"restaurantId"`
	UserId    uuid.UUID  `json:"userId"`
	FirstName string     `json:"frstName"`
	LastName  string     `json:"lastName"`
	Avatar    *string    `json:"avatar"`
	Rating    float64    `json:"rating"`
	Comment   *string    `json:"comment"`
	CreatedAt *time.Time `json:"createdAt"`
}
