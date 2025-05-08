package foodgormmysql

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *FoodLikeRepo) FindByUserIdAndFoodId(ctx context.Context, userId, foodId uuid.UUID) (*foodmodel.FoodLike, error) {
	var foodLike foodmodel.FoodLike
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("food_id = ? AND user_id = ?", foodId, userId).First(&foodLike).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, foodmodel.ErrFoodLikeNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &foodLike, nil
}
