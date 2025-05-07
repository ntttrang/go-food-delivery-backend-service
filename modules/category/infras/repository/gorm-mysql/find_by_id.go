package categorygormmysql

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *CategoryRepo) FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error) {
	var category categorymodel.Category
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return categorymodel.Category{}, categorymodel.ErrCategoryNotFound
		}
		return categorymodel.Category{}, errors.WithStack(err)
	}
	return category, nil
}
