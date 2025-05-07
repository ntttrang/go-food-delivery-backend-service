package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

// FindAll retrieves all active restaurants from the database
func (r *RestaurantRepo) FindAll(ctx context.Context) ([]restaurantmodel.Restaurant, error) {
	var restaurants []restaurantmodel.Restaurant
	
	if err := r.dbCtx.GetMainConnection().
		Where("status = ?", "ACTIVE").
		Find(&restaurants).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	
	return restaurants, nil
}
