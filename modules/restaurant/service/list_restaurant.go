package service

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListRestaurantRepo interface {
	List(ctx context.Context, req restaurantmodel.RestaurantListReq) ([]restaurantmodel.RestaurantSearchResDto, int64, error)
}

type ICategoryRepo interface {
}

type ListQueryHandler struct {
	restaurantRepo IListRestaurantRepo
	categoryRepo   ICategoryRepo
}

func NewListQueryHandler(restaurantRepo IListRestaurantRepo, categoryRepo ICategoryRepo) *ListQueryHandler {
	return &ListQueryHandler{
		restaurantRepo: restaurantRepo,
		categoryRepo:   categoryRepo,
	}
}

func (hdl *ListQueryHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantListReq) (restaurantmodel.RestaurantSearchRes, error) {
	restaurants, total, err := hdl.restaurantRepo.List(ctx, req)

	if err != nil {
		return restaurantmodel.RestaurantSearchRes{}, err
	}

	var resp restaurantmodel.RestaurantSearchRes
	resp.Items = restaurants
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return resp, nil
}
