package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantDeleteReq struct {
	Id uuid.UUID
}

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

func (hdl *DeleteCommandHandler) Execute(ctx context.Context, req RestaurantDeleteReq) error {
	existRestaurant, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if existRestaurant.Status == sharedModel.StatusDelete {
		return datatype.ErrNotFound.WithError(restaurantmodel.ErrRestaurantIsDeleted.Error())
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
