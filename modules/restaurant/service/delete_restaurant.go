package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IDeleteRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleteCommandHandler struct {
	repo IDeleteRepo
}

func NewDeleteCommandHandler(repo IDeleteRepo) *DeleteCommandHandler {
	return &DeleteCommandHandler{repo: repo}
}

func (hdl *DeleteCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantDeleteReq) error {
	existRestaurant, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		return err
	}

	if existRestaurant.Status == sharedModel.StatusDelete {
		return restaurantmodel.ErrRestaurantIsDeleted
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		return err
	}

	return nil
}
