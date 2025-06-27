package cartgormmysql

import (
	"context"

	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
)

func (r *CartRepo) Update(ctx context.Context, cart *cartmodel.Cart) error {
	db := r.dbCtx.GetMainConnection().Begin()

	if err := db.Save(cart).Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	return nil
}
