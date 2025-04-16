package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) List(ctx context.Context, req restaurantmodel.RestaurantListReq) ([]restaurantmodel.RestaurantSearchResDto, int64, error) {
	db := r.dbCtx.GetMainConnection().Table(restaurantmodel.Restaurant{}.TableName()).Select("id", "name", "addr", "logo", "shipping_fee_per_km", "status") // Use field name ( Struct) or gorm name is OK
	if req.OwnerId != nil {
		db = db.Where("owner_id = ?", req.OwnerId)
	}

	if req.CityId != nil {
		db = db.Where("city_id = ?", req.CityId)
	}
	db = db.Where("status = ?", usermodel.StatusActive)

	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}

	var result []restaurantmodel.RestaurantSearchResDto
	var total int64
	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&result).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return result, total, nil
}
