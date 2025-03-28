package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
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
	// 	return err
	// }

	existRestaurant, err := hdl.restaurantRepo.FindById(ctx, req.Id)

	if err != nil {
		return err
	}

	if existRestaurant.Status == sharedModel.StatusDelete {
		return restaurantmodel.ErrRestaurantIsDeleted
	}

	if err := hdl.restaurantRepo.Update(ctx, req); err != nil {
		return err
	}

	return nil
}
