package service

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListFavoritesRestaurantRepo interface {
	FindFavRestaurant(ctx context.Context, req restaurantmodel.FavoriteRestaurantListReq) ([]restaurantmodel.RestaurantSearchResDto, int64, error)
}

type ListFavoritesRestaurantQueryHandler struct {
	repo IListFavoritesRestaurantRepo
}

func NewGetFavoritesRestaurantQueryHandler(repo IListFavoritesRestaurantRepo) *ListFavoritesRestaurantQueryHandler {
	return &ListFavoritesRestaurantQueryHandler{repo: repo}
}

func (hdl *ListFavoritesRestaurantQueryHandler) Execute(ctx context.Context, req restaurantmodel.FavoriteRestaurantListReq) (restaurantmodel.RestaurantSearchRes, error) {
	restaurants, total, err := hdl.repo.FindFavRestaurant(ctx, req)

	if err != nil {
		return restaurantmodel.RestaurantSearchRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp restaurantmodel.RestaurantSearchRes
	resp.Items = restaurants
	if restaurants == nil {
		resp.Items = make([]restaurantmodel.RestaurantSearchResDto, 0)
	}

	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return resp, nil
}
