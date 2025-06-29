package service

import (
	"context"
	"encoding/json"
	"time"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type OrderSearchDto struct {
	UserID       *string `json:"userId" form:"userId"`
	RestaurantID *string `json:"restaurantId" form:"restaurantId"`
	Status       *string `json:"status" form:"status"`
	State        *string `json:"state" form:"state"`
}

type OrderListReq struct {
	OrderSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

type OrderListRes struct {
	Data   []OrderDto            `json:"data"`
	Paging sharedModel.PagingDto `json:"paging"`
}

type OrderDto struct {
	ID              string             `json:"id"`
	UserID          string             `json:"userId"`
	TotalPrice      float64            `json:"totalPrice"`
	ShipperID       *string            `json:"shipperId,omitempty"`
	Status          string             `json:"status"`
	State           string             `json:"state"`
	PaymentStatus   string             `json:"paymentStatus"`
	PaymentMethod   string             `json:"paymentMethod"`
	DeliveryAddress ordermodel.Address `json:"deliveryAddress"`
	RestaurantID    string             `json:"restaurantId"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

// Initialize service
type IListOrderRepo interface {
	List(ctx context.Context, req OrderListReq) ([]ordermodel.Order, []ordermodel.OrderTracking, int64, error)
}

type ListQueryHandler struct {
	repo IListOrderRepo
}

func NewListQueryHandler(repo IListOrderRepo) *ListQueryHandler {
	return &ListQueryHandler{repo: repo}
}

// Implement
func (hdl *ListQueryHandler) Execute(ctx context.Context, req OrderListReq) (*OrderListRes, error) {
	// Process paging
	req.PagingDto.Process()

	// Get orders from database
	orders, trackings, total, err := hdl.repo.List(ctx, req)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Create tracking map for quick lookup
	trackingMap := make(map[string]ordermodel.OrderTracking)
	for _, tracking := range trackings {
		trackingMap[tracking.OrderID] = tracking
	}

	// Map to response DTOs
	var orderDtos []OrderDto
	for _, order := range orders {
		tracking, exists := trackingMap[order.ID]
		if !exists {
			continue
		}

		var deliveryAddress ordermodel.Address
		json.Unmarshal(tracking.DeliveryAddress, &deliveryAddress)

		orderDtos = append(orderDtos, OrderDto{
			ID:              order.ID,
			UserID:          order.UserID,
			TotalPrice:      order.TotalPrice,
			ShipperID:       order.ShipperID,
			Status:          order.Status,
			State:           tracking.State,
			PaymentStatus:   tracking.PaymentStatus,
			PaymentMethod:   tracking.PaymentMethod,
			DeliveryAddress: deliveryAddress,
			RestaurantID:    tracking.RestaurantID,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
		})
	}

	return &OrderListRes{
		Data: orderDtos,
		Paging: sharedModel.PagingDto{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}, nil
}
