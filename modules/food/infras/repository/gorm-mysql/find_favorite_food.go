package foodgormmysql

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/pkg/errors"
)

func (r *FoodRepo) FindFavFood(ctx context.Context, req service.FavoriteFoodListReq) ([]foodmodel.FoodSearchResDto, int64, error) {
	db := r.dbCtx.GetMainConnection().Table("foods f").
		Select("f.id", "f.name", "f.description", "f.status", "f.created_at", "f.updated_at"). // Use field name ( Struct) or gorm name is OK
		Joins("INNER JOIN food_likes fl ON fl.food_id = f.id").
		Where("fl.user_id = ?", req.UserId)

	sortStr := "f.created_at DESC"
	var results []foodmodel.FoodSearchResDto
	var total int64
	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Scan(&results).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return results, total, nil
}
