package service

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListFavoritesFoodRepo interface {
	FindFavFood(ctx context.Context, req foodmodel.FavoriteFoodListReq) ([]foodmodel.FoodSearchResDto, int64, error)
}

type ListFavoritesFoodQueryHandler struct {
	repo IListFavoritesFoodRepo
}

func NewGetFavoritesFoodQueryHandler(repo IListFavoritesFoodRepo) *ListFavoritesFoodQueryHandler {
	return &ListFavoritesFoodQueryHandler{repo: repo}
}

func (hdl *ListFavoritesFoodQueryHandler) Execute(ctx context.Context, req foodmodel.FavoriteFoodListReq) (foodmodel.ListFoodRes, error) {
	foods, total, err := hdl.repo.FindFavFood(ctx, req)

	if err != nil {
		return foodmodel.ListFoodRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp foodmodel.ListFoodRes
	resp.Items = foods
	if foods == nil {
		resp.Items = make([]foodmodel.FoodSearchResDto, 0)
	}

	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return resp, nil
}
