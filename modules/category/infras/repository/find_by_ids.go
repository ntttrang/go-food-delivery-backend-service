package repository

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (r *CategoryRepo) FindByIds(ctx context.Context, ids []uuid.UUID) ([]categorymodel.Category, error) {
	var categories []categorymodel.Category
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id IN (?)", ids).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
