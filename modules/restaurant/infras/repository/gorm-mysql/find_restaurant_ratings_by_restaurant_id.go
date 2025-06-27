package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *RestaurantRatingRepo) FindByRestaurantIdOrUserId(ctx context.Context, req restaurantservice.RestaurantRatingListReq) ([]restaurantmodel.RestaurantRating, int64, error) {
	var restaurantRatings []restaurantmodel.RestaurantRating

	tx := r.dbCtx.GetMainConnection().
		Model(&restaurantmodel.RestaurantRating{}).
		Select("id, comment, point, created_at, user_id, restaurant_id").
		Where("status = ?", datatype.StatusActive)

	if req.RestaurantId != "" {
		tx = tx.Where("restaurant_id = ?", req.RestaurantId)
	}
	if req.UserId != "" {
		tx = tx.Where("user_id = ?", req.UserId)
	}

	sortStr := "created_at DESC"

	var total int64
	if err := tx.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&restaurantRatings).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return restaurantRatings, total, nil
}
