package cartgormmysql

import (
	"context"

	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
)

func (r *CartRepo) Insert(ctx context.Context, cart *cartmodel.Cart) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Create(cart).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
