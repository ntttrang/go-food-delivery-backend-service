package repository

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *CategoryRepo) FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error) {
	var category categorymodel.Category

	if err := r.db.Where("id = ?", id).Find(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return categorymodel.Category{}, categorymodel.ErrCategoryNotFound
		}
		return categorymodel.Category{}, errors.WithStack(err)
	}
	return category, nil
}
