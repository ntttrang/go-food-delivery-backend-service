package foodmodel

import (
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ListFoodReq struct {
	SearchFoodDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}
