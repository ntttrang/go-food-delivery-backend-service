package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) List(ctx context.Context, req restaurantservice.RestaurantListReq) ([]restaurantservice.RestaurantSearchResDto, int64, error) {
	db := r.dbCtx.GetMainConnection().Table(restaurantmodel.Restaurant{}.TableName()).Select("id", "name", "addr AS address", "logo", "shipping_fee_per_km", "status") // Use field name ( Struct) or gorm name is OK
	if req.OwnerId != nil {
		db = db.Where("owner_id = ?", req.OwnerId)
	}

	if req.CityId != nil {
		db = db.Where("city_id = ?", req.CityId)
	}
	db = db.Where("status = ?", datatype.StatusActive)

	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}

	var modelResult []restaurantservice.RestaurantSearchResDto
	var total int64
	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&modelResult).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	// Convert model DTOs to service DTOs
	result := make([]restaurantservice.RestaurantSearchResDto, len(modelResult))
	for i, item := range modelResult {
		result[i] = restaurantservice.RestaurantSearchResDto{
			ID:               item.ID,
			Name:             item.Name,
			Address:          item.Address,
			Logo:             item.Logo,
			Cover:            item.Cover,
			ShippingFeePerKm: item.ShippingFeePerKm,
			Status:           item.Status,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		}
	}

	return result, total, nil
}
