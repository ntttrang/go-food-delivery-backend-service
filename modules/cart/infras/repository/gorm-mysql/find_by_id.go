package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *CartRepo) FindById(ctx context.Context, id uuid.UUID) ([]cartmodel.Cart, error) {
	var carts []cartmodel.Cart
	db := r.dbCtx.GetMainConnection()

	if err := db.Where("id = ?", id).Find(&carts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cartmodel.ErrCartNotFound
		}
		return nil, errors.WithStack(err)
	}

	return carts, nil
}
