package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
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
		return err
	}

	if err := hdl.repo.Delete(ctx, req.RestaurantID, req.UserID); err != nil {
		return err
	}

	return nil
}
