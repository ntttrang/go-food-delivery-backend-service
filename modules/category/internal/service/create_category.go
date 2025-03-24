package service

import (
	"context"

	"github.com/google/uuid"
	categoryModel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (s *CategoryService) CreateCategory(ctx context.Context, data *categoryModel.CategoryInsertDto) error {
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
