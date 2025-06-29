package ordergormmysql

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *OrderRepo) List(ctx context.Context, req service.OrderListReq) ([]ordermodel.Order, []ordermodel.OrderTracking, int64, error) {
	db := r.dbCtx.GetMainConnection()
	var orders []ordermodel.Order
	var trackings []ordermodel.OrderTracking
	var total int64

	// Build query for orders
	orderQuery := db.Model(&ordermodel.Order{})

	// Apply filters
	if req.UserID != nil {
		orderQuery = orderQuery.Where("user_id = ?", *req.UserID)
	}

	if req.Status != nil {
		orderQuery = orderQuery.Where("status = ?", *req.Status)
	} else {
		// By default, only return active orders
		orderQuery = orderQuery.Where("status = ?", datatype.StatusActive)
	}

	// Count total records
	if err := orderQuery.Count(&total).Error; err != nil {
		return nil, nil, 0, errors.WithStack(err)
	}

	// Apply pagination
	orderQuery = orderQuery.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)

	// Apply sorting
	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}
	orderQuery = orderQuery.Order(sortStr)

	// Execute query
	if err := orderQuery.Find(&orders).Error; err != nil {
		return nil, nil, 0, errors.WithStack(err)
	}

	// If no orders found, return empty result
	if len(orders) == 0 {
		return orders, trackings, total, nil
	}

	// Get order IDs
	var orderIDs []string
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	// Build query for order trackings
	trackingQuery := db.Model(&ordermodel.OrderTracking{}).Where("order_id IN ?", orderIDs)

	// Apply restaurant filter if provided
	if req.RestaurantID != nil {
		trackingQuery = trackingQuery.Where("restaurant_id = ?", *req.RestaurantID)
	}

	// Apply state filter if provided
	if req.State != nil {
		trackingQuery = trackingQuery.Where("state = ?", *req.State)
	}

	// Execute tracking query
	if err := trackingQuery.Find(&trackings).Error; err != nil {
		return nil, nil, 0, errors.WithStack(err)
	}

	return orders, trackings, total, nil
}
