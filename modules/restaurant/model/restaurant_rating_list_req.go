package restaurantmodel

import (
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantRatingListReq struct {
	RestaurantId string `json:"restaurantId" form:"restaurantId"`
	UserId       string `json:"userId" form:"userId"`
	sharedModel.PagingDto
}

func (r *RestaurantRatingListReq) Validate() error {
	if r.RestaurantId == "" && r.UserId == "" {
		return ErrFieldRequired
	}

	return nil
}
