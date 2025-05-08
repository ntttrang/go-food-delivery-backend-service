package service

import (
	"context"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type UserListReq struct {
	UserSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

type UserSearchDto struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

type UserListRes struct {
	Items      []UserSearchResDto    `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

// Initilize service
type IListUserRepo interface {
	FindUsers(ctx context.Context, req UserListReq) ([]UserSearchResDto, int64, error)
}

type ListQueryHandler struct {
	userRepo IListUserRepo
}

func NewListQueryHandler(userRepo IListUserRepo) *ListQueryHandler {
	return &ListQueryHandler{
		userRepo: userRepo,
	}
}

// Implement
func (hdl *ListQueryHandler) Execute(ctx context.Context, req UserListReq) (UserListRes, error) {
	users, total, err := hdl.userRepo.FindUsers(ctx, req)

	if err != nil {
		return UserListRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp UserListRes
	resp.Items = users
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return resp, nil
}
