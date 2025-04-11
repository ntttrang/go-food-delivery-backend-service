package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *RestaurantRepo) FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error) {
	var restaurant restaurantmodel.Restaurant

	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id = ?", id.String()).First(&restaurant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, restaurantmodel.ErrRestaurantNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &restaurant, nil
}
