package restaurantmodel

import (
	"time"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantRatingListRes struct {
	Items      []RestaurantRatingListDto `json:"items"`
	Pagination sharedModel.PagingDto     `json:"pagination"`
}

type RestaurantRatingListDto struct {
	Id           uuid.UUID  `json:"id"`
	RestaurantId uuid.UUID  `json:"restaurantId"`
	UserId       uuid.UUID  `json:"userId"`
	FirstName    string     `json:"frstName"`
	LastName     string     `json:"lastName"`
	Avatar       *string    `json:"avatar"`
	Rating       float64    `json:"rating"`
	Comment      *string    `json:"comment"`
	CreatedAt    *time.Time `json:"createdAt"`
}
