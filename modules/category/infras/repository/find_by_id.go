package repository

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (r *CategoryRepository) FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error) {
	var category categorymodel.Category

	if err := r.db.Where("id = ?", id).Find(&category).Error; err != nil {
		return categorymodel.Category{}, err
	}
	return category, nil
}
