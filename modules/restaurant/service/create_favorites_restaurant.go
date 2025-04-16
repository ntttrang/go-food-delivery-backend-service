package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IRestaurantLikeRepo interface {
	FindByUserIdAndRestaurantId(ctx context.Context, userId, restaurantId uuid.UUID) (*restaurantmodel.RestaurantLike, error)
	Insert(ctx context.Context, req restaurantmodel.RestaurantLike) error
	Delete(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) error
}

type AddFavoritesCommandHandler struct {
	repo IRestaurantLikeRepo
}

func NewAddFavoritesCommandHandler(repo IRestaurantLikeRepo) *AddFavoritesCommandHandler {
	return &AddFavoritesCommandHandler{
		repo: repo,
	}
}

func (hdl *AddFavoritesCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantLike) (*string, error) {
	if err := req.Validate(); err != nil {
		return nil, datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	restaurantLike, err := hdl.repo.FindByUserIdAndRestaurantId(ctx, req.UserID, req.RestaurantID)

	if err != nil {
		if !errors.Is(err, restaurantmodel.ErrRestaurantLikeNotFound) {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
	}

	flag := ""
	if restaurantLike == nil {
		if err := hdl.repo.Insert(ctx, req); err != nil {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
		flag = "Add to My favorite restaurant"
	} else {
		if err := hdl.repo.Delete(ctx, req.RestaurantID, req.UserID); err != nil {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
		flag = "Remove out My favorite restaurant"
	}

	return &flag, nil
}
