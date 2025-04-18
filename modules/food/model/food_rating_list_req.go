package foodmodel

import (
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FoodRatingListReq struct {
	FoodId string `json:"foodId" form:"foodId"`
	UserId string `json:"userId" form:"userId"`
	sharedModel.PagingDto
}

func (r *FoodRatingListReq) Validate() error {
	if r.FoodId == "" && r.UserId == "" {
		return ErrFieldRequired
	}

	return nil
}
