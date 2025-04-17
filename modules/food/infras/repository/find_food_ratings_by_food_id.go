package repository

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

func (r *FoodRatingRepo) FindByFoodIdOrUserId(ctx context.Context, req foodmodel.FoodRatingListReq) ([]foodmodel.FoodRatings, int64, error) {
	var foodRatings []foodmodel.FoodRatings

	tx := r.dbCtx.GetMainConnection().
		Model(&foodmodel.FoodRatings{}).
		Select("id, comment, point, created_at, user_id, food_id").
		Where("status = ?", sharedModel.StatusActive)

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
