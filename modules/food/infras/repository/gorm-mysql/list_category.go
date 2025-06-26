package foodgormmysql

import (
	"context"

	"github.com/pkg/errors"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (r *FoodRepo) ListFoods(ctx context.Context, req service.ListFoodReq) ([]foodmodel.Food, int64, error) {

	var result []foodmodel.Food
	var total int64

	db := r.dbCtx.GetMainConnection().Table(foodmodel.Food{}.TableName())

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Description != "" {
		db = db.Where("description LIKE ?", "%"+req.Description+"%")
	}

	db = db.Where("status in (?)", []string{string(datatype.StatusActive)})

	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}
	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&result).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return result, total, nil
}
