package service

import (
	"context"

	"github.com/google/uuid"
	categoryModel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (s *CategoryService) BulkInsert(ctx context.Context, datas []categoryModel.CategoryInsertDto) ([]uuid.UUID, error) {

	var categories []categoryModel.Category
	var ids []uuid.UUID
	for _, data := range datas {
		if err := data.Validate(); err != nil {
			return nil, err
		}

		category := data.ConvertToCategory()
		category.Id, _ = uuid.NewV7()
		category.Status = sharedModel.StatusActive // Always set Active Status when insert

		categories = append(categories, *category)
		ids = append(ids, category.Id)
	}

	if err := s.catRepo.BulkInsert(ctx, categories); err != nil {
		return nil, err
	}

	// set data to response
	return ids, nil
}
