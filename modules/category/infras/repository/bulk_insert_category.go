package repository

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (r *CategoryRepository) BulkInsert(ctx context.Context, data []categorymodel.Category) error {
	if err := r.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
