package cartgormmysql

import (
	"context"

	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	"github.com/pkg/errors"
)

func (r *CartRepo) ListByUserId(ctx context.Context, req service.CartListReq) ([]cartmodel.Cart, int64, error) {
	var carts []cartmodel.Cart
	var total int64

	db := r.dbCtx.GetMainConnection().Table(cartmodel.Cart{}.TableName())

	if req.UserID != nil {
		db = db.Where("user_id = ?", req.UserID)
	}

	// By default, only return active carts
	db = db.Where("status != ?", cartmodel.CartStatusProcessed)

	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}

	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&carts).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return carts, total, nil
}
