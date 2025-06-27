package cartgormmysql

import (
	"context"

	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *CartRepo) ListByUserId(ctx context.Context, req service.CartListReq) ([]cartmodel.Cart, error) {
	var carts []cartmodel.Cart

	db := r.dbCtx.GetMainConnection().Table(cartmodel.Cart{}.TableName()).Select("id", "user_id", "restaurant_id", "COUNT(*) AS item_quantity", "dropoff_lat", "dropoff_lng")

	if req.UserID != nil {
		db = db.Where("user_id = ?", req.UserID)
	}

	// By default, only return active carts
	db = db.Where("status != ?", datatype.CartStatusProcessed).Group("id,user_id, restaurant_id, dropoff_lat, dropoff_lng")

	// TODO: Sort newest cart
	// sortStr := "created_at DESC"
	// if req.SortBy != "" {
	// 	sortStr = req.SortBy + " " + req.Direction
	// }

	if err := db.Find(&carts).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return carts, nil
}
