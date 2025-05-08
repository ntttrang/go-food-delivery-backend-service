package restaurantgormmysql

import (
	"context"

	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) FindFavRestaurant(ctx context.Context, req restaurantservice.FavoriteRestaurantListReq) ([]restaurantservice.RestaurantSearchResDto, int64, error) {
	db := r.dbCtx.GetMainConnection().Table("restaurants r").
		Select("r.id", "r.name", "r.addr", "r.logo", "r.shipping_fee_per_km", "r.status"). // Use field name ( Struct) or gorm name is OK
		Joins("INNER JOIN restaurant_likes rl ON rl.restaurant_id = r.id").
		Where("rl.user_id = ?", req.UserId)

	sortStr := "r.created_at DESC"
	var results []restaurantservice.RestaurantSearchResDto
	var total int64
	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Scan(&results).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return results, total, nil
}
