package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type MenuItemCreateReq struct {
	RestaurantId uuid.UUID `json:"restaurantId" form:"restaurantId"`
	FoodId       uuid.UUID `json:"foodId" form:"foodId"`
	CategoryId   uuid.UUID `json:"categoryId" form:"categoryId"`
}

func (r MenuItemCreateReq) ConvertToRestaurantFood() *restaurantmodel.RestaurantFood {
	return &restaurantmodel.RestaurantFood{
		RestaurantId: r.RestaurantId,
		FoodId:       r.FoodId,
	}
}

func (r MenuItemCreateReq) Validate() error {
	if r.RestaurantId.String() == "" {
		return restaurantmodel.ErrRestaurantIdRequired
	}
	if r.FoodId.String() == "" {
		return restaurantmodel.ErrFoodIdRequired
	}
	if r.CategoryId.String() == "" {
		return restaurantmodel.ErrCategoryIdRequired
	}
	return nil
}

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

func (s *CreateMenuItemCommandHandler) Execute(ctx context.Context, req *MenuItemCreateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}
	categoryId := req.CategoryId
	restaurantFood := req.ConvertToRestaurantFood()
	restaurantFood.Status = string(datatype.StatusActive) // Always set Active Status when insert
	now := time.Now().UTC()
	restaurantFood.CreatedAt = &now
	restaurantFood.UpdatedAt = &now

	if err := s.menuItemRepo.Insert(ctx, restaurantFood, categoryId); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
