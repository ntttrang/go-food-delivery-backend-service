package service

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

type IRestaurantLikeRepo interface {
	Insert(ctx context.Context, req restaurantmodel.RestaurantLike) error
}

type AddFavoritesCommandHandler struct {
	repo IRestaurantLikeRepo
}

func NewAddFavoritesCommandHandler(repo IRestaurantLikeRepo) *AddFavoritesCommandHandler {
	return &AddFavoritesCommandHandler{
		repo: repo,
	}
}

func (hdl *AddFavoritesCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantLike) error {
	if err := req.Validate(); err != nil {
		return err
	}

	if err := hdl.repo.Insert(ctx, req); err != nil {
		return err
	}

	return nil
}
