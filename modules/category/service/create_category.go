package service

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ICreateRepo interface {
	Insert(ctx context.Context, data *categorymodel.Category) error
}

type CreateCommandHandler struct {
	catRepo ICreateRepo
}

func NewCreateCommandHandler(catRepo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{catRepo: catRepo}
}

func (s *CreateCommandHandler) Execute(ctx context.Context, data *categorymodel.CategoryInsertDto) error {
	if err := data.Validate(); err != nil {
		return err
	}

	category := data.ConvertToCategory()
	category.Id, _ = uuid.NewV7()
	category.Status = sharedModel.StatusActive // Always set Active Status when insert

	if err := s.catRepo.Insert(ctx, category); err != nil {
		return err
	}

	// set data to response
	data.Id = category.Id

	return nil
}
