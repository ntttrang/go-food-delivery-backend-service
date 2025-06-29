package service

import (
	"context"
	"encoding/json"
	"time"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type OrderDetailReq struct {
	ID string `json:"-"`
}

type OrderDetailRes struct {
	ID              string               `json:"id"`
	UserID          string               `json:"userId"`
	TotalPrice      float64              `json:"totalPrice"`
	ShipperID       *string              `json:"shipperId,omitempty"`
	Status          string               `json:"status"`
	State           string               `json:"state"`
	PaymentStatus   string               `json:"paymentStatus"`
	PaymentMethod   string               `json:"paymentMethod"`
	DeliveryAddress ordermodel.Address   `json:"deliveryAddress"`
	DeliveryFee     float64              `json:"deliveryFee"`
	EstimatedTime   int                  `json:"estimatedTime"`
	DeliveryTime    int                  `json:"deliveryTime"`
	RestaurantID    string               `json:"restaurantId"`
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
	OrderDetails    []OrderDetailItemDto `json:"orderDetails"`
}

type OrderDetailItemDto struct {
	ID         string        `json:"id"`
	OrderID    string        `json:"orderId"`
	FoodOrigin FoodOriginDto `json:"foodOrigin"`
	Price      float64       `json:"price"`
	Quantity   int           `json:"quantity"`
	Discount   float64       `json:"discount"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

// Initialize service
type IGetOrderDetailRepo interface {
	FindById(ctx context.Context, id string) (*ordermodel.Order, *ordermodel.OrderTracking, []ordermodel.OrderDetail, error)
}

type GetDetailQueryHandler struct {
	repo IGetOrderDetailRepo
}

func NewGetDetailQueryHandler(repo IGetOrderDetailRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

// Implement
func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req OrderDetailReq) (*OrderDetailRes, error) {
	if req.ID == "" {
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrOrderIdRequired).WithDebug(ordermodel.ErrOrderIdRequired.Error())
	}

	// Get order from database
	order, tracking, details, err := hdl.repo.FindById(ctx, req.ID)
	if err != nil {
		if err == ordermodel.ErrOrderNotFound {
			return nil, datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Map order details to DTOs
	var orderDetailDtos []OrderDetailItemDto
	for _, detail := range details {
		var foodOrigin FoodOriginDto
		json.Unmarshal(detail.FoodOrigin, &foodOrigin)
		orderDetailDtos = append(orderDetailDtos, OrderDetailItemDto{
			ID:         detail.ID,
			OrderID:    detail.OrderID,
			FoodOrigin: foodOrigin,
			Price:      detail.Price,
			Quantity:   detail.Quantity,
			Discount:   detail.Discount,
			Status:     detail.Status,
			CreatedAt:  detail.CreatedAt,
			UpdatedAt:  detail.UpdatedAt,
		})
	}

	// Create response
	var deliveryAddress ordermodel.Address
	json.Unmarshal(tracking.DeliveryAddress, &deliveryAddress)
	return &OrderDetailRes{
		ID:              order.ID,
		UserID:          order.UserID,
		TotalPrice:      order.TotalPrice,
		ShipperID:       order.ShipperID,
		Status:          order.Status,
		State:           tracking.State,
		PaymentStatus:   tracking.PaymentStatus,
		PaymentMethod:   tracking.PaymentMethod,
		DeliveryAddress: deliveryAddress,
		DeliveryFee:     tracking.DeliveryFee,
		EstimatedTime:   tracking.EstimatedTime,
		DeliveryTime:    tracking.DeliveryTime,
		RestaurantID:    tracking.RestaurantID,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		OrderDetails:    orderDetailDtos,
	}, nil
}
