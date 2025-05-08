package foodgormmysql

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *FoodRepo) FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error) {
	var food foodmodel.Food
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id = ?", id).First(&food).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return foodmodel.Food{}, foodmodel.ErrFoodNotFound
		}
		return foodmodel.Food{}, errors.WithStack(err)
	}
	return food, nil
}
