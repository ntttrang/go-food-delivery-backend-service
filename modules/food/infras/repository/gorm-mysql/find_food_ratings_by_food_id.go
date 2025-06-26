package foodgormmysql

import (
	"context"

	"github.com/pkg/errors"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (r *FoodRatingRepo) FindByFoodIdOrUserId(ctx context.Context, req service.FoodRatingListReq) ([]foodmodel.FoodRatings, int64, error) {
	var foodRatings []foodmodel.FoodRatings

	tx := r.dbCtx.GetMainConnection().
		Model(&foodmodel.FoodRatings{}).
		Select("id, comment, point, created_at, user_id, food_id").
		Where("status = ?", string(datatype.StatusActive))

	if req.FoodId != "" {
		tx = tx.Where("food_id = ?", req.FoodId)
	}
	if req.UserId != "" {
		tx = tx.Where("user_id = ?", req.UserId)
	}

	sortStr := "created_at DESC"

	var total int64
	if err := tx.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&foodRatings).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return foodRatings, total, nil
}
