package usermodel

import sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"

type UserListReq struct {
	UserSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}
