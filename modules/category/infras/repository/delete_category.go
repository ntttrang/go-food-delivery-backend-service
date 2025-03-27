package repository

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.Table(categorymodel.Category{}.TableName()).
		Where("id = ?", id).
		Update("status", sharedModel.StatusDelete).
		Error; err != nil {
		return err
	}
	return nil
}
