package ordergormmysql

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *OrderRepo) FindById(ctx context.Context, id string) (*ordermodel.Order, *ordermodel.OrderTracking, []ordermodel.OrderDetail, error) {
	db := r.dbCtx.GetMainConnection()
	
	// Find order
	var order ordermodel.Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ordermodel.ErrOrderNotFound
		}
		return nil, nil, nil, errors.WithStack(err)
	}

	// Find order tracking
	var tracking ordermodel.OrderTracking
	if err := db.Where("order_id = ?", id).First(&tracking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Order exists but tracking doesn't - this is an inconsistent state
			return &order, nil, nil, errors.WithStack(err)
		}
		return &order, nil, nil, errors.WithStack(err)
	}

	// Find order details
	var details []ordermodel.OrderDetail
	if err := db.Where("order_id = ?", id).Find(&details).Error; err != nil {
		return &order, &tracking, nil, errors.WithStack(err)
	}

	return &order, &tracking, details, nil
}
