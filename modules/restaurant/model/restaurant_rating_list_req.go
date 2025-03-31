package restaurantmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantRatingListReq struct {
	RestaurantId string `json:"restaurantId" form:"restaurantId"`
	sharedModel.PagingDto
}

type RestaurantRatingSearchDto struct {
	RestaurantId uuid.UUID `json:"restaurantId" form:"restaurantId"`
}
