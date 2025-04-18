package restaurantmodel

import sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"

type RestaurantListReq struct {
	RestaurantSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}
