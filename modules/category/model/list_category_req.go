package categorymodel

import (
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ListCategoryReq struct {
	SearchCategoryDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}
