package service

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
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
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	if err := hdl.repo.Insert(ctx, req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
