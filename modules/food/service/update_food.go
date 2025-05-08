package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type FoodUpdateReq struct {
	// Use pointer to accept empty string
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Status       *string `json:"status"`
	RestaurantId *string `json:"restaurantId"` // Can be empty or missing if data type = string. Otherwise, uuid.UUID isn't
	CategoryId   *string `json:"categoryId"`   // Can be empty or missing if data type = string. Otherwise, uuid.UUID isn't

	Id uuid.UUID `json:"-"`
}

func (FoodUpdateReq) TableName() string {
	return foodmodel.Food{}.TableName()
}

func (c FoodUpdateReq) Validate() error {
	if c.Status != nil && *c.Status != sharedModel.StatusActive && *c.Status != sharedModel.StatusDelete && *c.Status != sharedModel.StatusInactive {
		return foodmodel.ErrFoodStatusInvalid
	}
	return nil
}

// Initilize service
type IUpdateByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
	Update(ctx context.Context, id uuid.UUID, dto FoodUpdateReq) error
}

type UpdateCommandHandler struct {
	repo IUpdateByIdRepo
}

func NewUpdateCommandHandler(repo IUpdateByIdRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: repo}
}

// Implement
func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req FoodUpdateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	category, err := hdl.repo.FindById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodNotFound) {
			return datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodNotFound.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if category.Status == sharedModel.StatusDelete {
		return datatype.ErrDeleted.WithError(foodmodel.ErrFoodIsDeleted.Error())
	}

	if err := hdl.repo.Update(ctx, req.Id, req); err != nil {
		return err
	}

	return nil
}
