package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (r *RestaurantRepo) FindRestaurantByIds(ctx context.Context, ids []uuid.UUID) ([]restaurantservice.FindRestaurantByIdsDto, error) {
	var restaurants []restaurantservice.FindRestaurantByIdsDto
	db := r.dbCtx.GetMainConnection()
	if len(ids) > 0 {
		if err := db.Raw(`SELECT r.id, 
				r.name, 
				r.addr, 
				r.city_id,
				r.Lat,
				r.Lng,
				r.Cover,
				r.Logo,
				r.shipping_fee_per_km, 
				COUNT(rr.comment) AS comment_qty,
				AVG(rr.point) AS avg_point,
				COUNT(rl.restaurant_id) AS likes_qty
			FROM restaurants r
			LEFT JOIN restaurant_ratings rr ON r.id = rr.restaurant_id
			LEFT JOIN restaurant_likes rl ON r.id = rl.restaurant_id
			WHERE r.id IN (?) AND r.status = ?
			GROUP BY r.id`, ids, sharedModel.StatusActive).
			Find(&restaurants).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Raw(`SELECT r.id, 
				r.name, 
				r.addr, 
				r.city_id,
				r.Lat,
				r.Lng,
				r.Cover,
				r.Logo,
				r.shipping_fee_per_km, 
				COUNT(rr.comment) AS comment_qty,
				AVG(rr.point) AS avg_point,
				COUNT(rl.restaurant_id) AS likes_qty
			FROM restaurants r
			LEFT JOIN restaurant_ratings rr ON r.id = rr.restaurant_id
			LEFT JOIN restaurant_likes rl ON r.id = rl.restaurant_id
			WHERE  r.status = ?
			GROUP BY r.id`, sharedModel.StatusActive).
			Find(&restaurants).Error; err != nil {
			return nil, err
		}
	}

	return restaurants, nil
}
