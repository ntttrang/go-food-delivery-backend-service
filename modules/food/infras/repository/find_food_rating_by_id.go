package repository

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *FoodRatingRepo) FindById(ctx context.Context, commentId uuid.UUID) (*foodmodel.FoodRatings, error) {
	var foodRating foodmodel.FoodRatings
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id = ?", commentId).First(&foodRating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, foodmodel.ErrFoodNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &foodRating, nil
}
