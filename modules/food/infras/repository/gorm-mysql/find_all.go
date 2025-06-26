package foodgormmysql

import (
	"context"

	"github.com/pkg/errors"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (r *FoodRepo) FindAll(ctx context.Context) ([]foodmodel.Food, error) {
	var result []foodmodel.Food

	db := r.dbCtx.GetMainConnection().Table(foodmodel.Food{}.TableName())

	// Only get active foods
	db = db.Where("status = ?", string(datatype.StatusActive))

	if err := db.Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
