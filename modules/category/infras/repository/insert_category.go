package repository

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/model"
)

func (repo *CategoryRepository) Insert(ctx context.Context, data *categorymodel.Category) error {
	if err := repo.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepository) BulkInsert(ctx context.Context, data []categorymodel.Category) error {
	if err := repo.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
