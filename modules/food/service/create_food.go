package service

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ICreateRepo interface {
	Insert(ctx context.Context, data *foodmodel.Food) error
}

type CreateCommandHandler struct {
	repo ICreateRepo
}

func NewCreateCommandHandler(repo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{repo: repo}
}

func (s *CreateCommandHandler) Execute(ctx context.Context, data *foodmodel.FoodInsertDto) error {
	if err := data.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	food := data.ConvertToFood()
	food.Id, _ = uuid.NewV7()
	food.Status = sharedModel.StatusActive // Always set Active Status when insert

	if err := s.repo.Insert(ctx, food); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	data.Id = food.Id

	return nil
}
