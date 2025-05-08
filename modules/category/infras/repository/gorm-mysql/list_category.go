package categorygormmysql

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

func (r *CategoryRepo) ListCategories(ctx context.Context, req service.ListCategoryReq) ([]categorymodel.Category, int64, error) {

	var categories []categorymodel.Category
	var total int64

	db := r.dbCtx.GetMainConnection().Table(categorymodel.Category{}.TableName())

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Description != "" {
		db = db.Where("description LIKE ?", "%"+req.Description+"%")
	}
	db = db.Where("status = ?", sharedModel.StatusActive)

	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}

	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&categories).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return categories, total, nil
}
