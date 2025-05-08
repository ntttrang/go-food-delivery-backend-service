package foodgormmysql

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

func (r *FoodRepo) FindAll(ctx context.Context) ([]foodmodel.Food, error) {
	var result []foodmodel.Food

	db := r.dbCtx.GetMainConnection().Table(foodmodel.Food{}.TableName())

	// Only get active foods
	db = db.Where("status = ?", sharedModel.StatusActive)

	if err := db.Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
