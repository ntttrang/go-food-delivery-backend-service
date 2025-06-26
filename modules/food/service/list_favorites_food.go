package service

import (
	"context"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type FavoriteFoodListReq struct {
	UserId uuid.UUID `json:"-" form:"-"`
	sharedModel.PagingDto
}

// Initilize service
type IListFavoritesFoodRepo interface {
	FindFavFood(ctx context.Context, req FavoriteFoodListReq) ([]foodmodel.FoodSearchResDto, int64, error)
}

type ListFavoritesFoodQueryHandler struct {
	repo IListFavoritesFoodRepo
}

func NewGetFavoritesFoodQueryHandler(repo IListFavoritesFoodRepo) *ListFavoritesFoodQueryHandler {
	return &ListFavoritesFoodQueryHandler{repo: repo}
}

// Implement
func (hdl *ListFavoritesFoodQueryHandler) Execute(ctx context.Context, req FavoriteFoodListReq) (*ListFoodRes, error) {
	foods, total, err := hdl.repo.FindFavFood(ctx, req)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp ListFoodRes
	resp.Items = foods
	if foods == nil {
		resp.Items = make([]foodmodel.FoodSearchResDto, 0)
	}

	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
