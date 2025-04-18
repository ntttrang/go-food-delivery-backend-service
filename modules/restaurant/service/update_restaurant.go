package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IUpdateRestaurantRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error)
	Update(ctx context.Context, req restaurantmodel.RestaurantUpdateReq) error
}

type UpdateRestaurantCommandHandler struct {
	restaurantRepo IUpdateRestaurantRepo
}

func NewUpdateRestaurantCommandHandler(restaurantRepo IUpdateRestaurantRepo) *UpdateRestaurantCommandHandler {
	return &UpdateRestaurantCommandHandler{restaurantRepo: restaurantRepo}
}

func (hdl *UpdateRestaurantCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantUpdateReq) error {
	// if err := req.Dto.Validate(); err != nil {
	// 	return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	// }

	existRestaurant, err := hdl.restaurantRepo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if existRestaurant.Status == sharedModel.StatusDelete {
		return datatype.ErrNotFound.WithError(restaurantmodel.ErrRestaurantIsDeleted.Error())
	}

	if err := hdl.restaurantRepo.Update(ctx, req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
