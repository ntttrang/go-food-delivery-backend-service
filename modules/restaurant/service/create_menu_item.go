package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ICreateMenuItemRepo interface {
	Insert(ctx context.Context, restaurant *restaurantmodel.RestaurantFood, categoryId uuid.UUID) error
}

type CreateMenuItemCommandHandler struct {
	menuItemRepo ICreateMenuItemRepo
}

func NewCreateMenuItemCommandHandler(menuItemRepo ICreateMenuItemRepo) *CreateMenuItemCommandHandler {
	return &CreateMenuItemCommandHandler{
		menuItemRepo: menuItemRepo,
	}
}

func (s *CreateMenuItemCommandHandler) Execute(ctx context.Context, req *restaurantmodel.MenuItemCreateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}
	categoryId := req.CategoryId
	restaurantFood := req.ConvertToRestaurantFood()
	restaurantFood.Status = sharedModel.StatusActive // Always set Active Status when insert
	now := time.Now().UTC()
	restaurantFood.CreatedAt = &now
	restaurantFood.UpdatedAt = &now

	if err := s.menuItemRepo.Insert(ctx, restaurantFood, categoryId); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
