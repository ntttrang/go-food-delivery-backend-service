package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IDeleteRestaurantLikeRepo interface {
	FindById(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) (*restaurantmodel.RestaurantLike, error)
	Delete(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) error
}

type DeleteRestaurantLikeCommandHandler struct {
	repo IDeleteRestaurantLikeRepo
}

func NewDeleteRestaurantLikeCommandHandler(repo IDeleteRestaurantLikeRepo) *DeleteRestaurantLikeCommandHandler {
	return &DeleteRestaurantLikeCommandHandler{repo: repo}
}

func (hdl *DeleteRestaurantLikeCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantLike) error {
	_, err := hdl.repo.FindById(ctx, req.RestaurantID, req.UserID)

	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if err := hdl.repo.Delete(ctx, req.RestaurantID, req.UserID); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
