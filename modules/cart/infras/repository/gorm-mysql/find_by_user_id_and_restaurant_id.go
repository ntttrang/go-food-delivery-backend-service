package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *CartRepo) FindByUserIdAndRestaurantId(ctx context.Context, userId, restaurantId uuid.UUID) (*cartmodel.Cart, error) {
	var cart cartmodel.Cart
	db := r.dbCtx.GetMainConnection()

	if err := db.Where("user_id = ? AND restaurant_id = ?", userId, restaurantId).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cartmodel.ErrCartNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &cart, nil
}
