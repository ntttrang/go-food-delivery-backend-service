package categorymodel

import (
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ListCategoryReq struct {
	SearchCategoryDto
	Paging sharedModel.PagingDto `json:"paging"`
}
