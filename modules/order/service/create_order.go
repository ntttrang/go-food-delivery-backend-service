package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type OrderCreateDto struct {
	UserID          string                 `json:"-"` // Get from token
	TotalPrice      float64                `json:"totalPrice"`
	RestaurantID    string                 `json:"restaurantId"`
	DeliveryAddress *ordermodel.Address    `json:"deliveryAddress"`
	PaymentMethod   string                 `json:"paymentMethod"`
	OrderDetails    []OrderDetailCreateDto `json:"orderDetails"`
}

type OrderDetailCreateDto struct {
	FoodOrigin *FoodOriginDto `json:"foodOrigin"`
	Price      float64        `json:"price"`
	Quantity   int            `json:"quantity"`
	Discount   float64        `json:"discount"`
}

type FoodOriginDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func (o *OrderCreateDto) Validate() error {
	if o.UserID == "" {
		return ordermodel.ErrUserIdRequired
	}

	if o.TotalPrice <= 0 {
		return ordermodel.ErrTotalPriceInvalid
	}

	if o.RestaurantID == "" {
		return ordermodel.ErrRestaurantRequired
	}

	return nil
}

// Initialize service
type ICreateOrderRepository interface {
	Insert(ctx context.Context, order *ordermodel.Order, orderTracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error
}

type CreateCommandHandler struct {
	repo ICreateOrderRepository
}

func NewCreateCommandHandler(repo ICreateOrderRepository) *CreateCommandHandler {
	return &CreateCommandHandler{repo: repo}
}

// Implement
func (s *CreateCommandHandler) Execute(ctx context.Context, data *OrderCreateDto) (string, error) {
	if err := data.Validate(); err != nil {
		return "", datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Generate new UUID for order
	orderId := uuid.New().String()

	// Create order
	order := &ordermodel.Order{
		ID:         orderId,
		UserID:     data.UserID,
		TotalPrice: data.TotalPrice,
		Status:     sharedModel.StatusActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Create order tracking
	var addressJson []byte
	if data.DeliveryAddress != nil {
		addressJson, _ = json.Marshal(data.DeliveryAddress)
	}
	orderTracking := &ordermodel.OrderTracking{
		ID:              uuid.New().String(),
		OrderID:         orderId,
		State:           "waiting_for_shipper",
		PaymentStatus:   "pending",
		PaymentMethod:   data.PaymentMethod,
		DeliveryAddress: addressJson,
		RestaurantID:    data.RestaurantID,
		Status:          sharedModel.StatusActive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create order details
	var orderDetails []ordermodel.OrderDetail
	for _, detail := range data.OrderDetails {
		var foodOriginJson []byte
		if detail.FoodOrigin != nil {
			foodOriginJson, _ = json.Marshal(detail.FoodOrigin)
		}
		orderDetails = append(orderDetails, ordermodel.OrderDetail{
			ID:         uuid.New().String(),
			OrderID:    orderId,
			FoodOrigin: foodOriginJson,
			Price:      detail.Price,
			Quantity:   detail.Quantity,
			Discount:   detail.Discount,
			Status:     sharedModel.StatusActive,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	// Insert to database
	if err := s.repo.Insert(ctx, order, orderTracking, orderDetails); err != nil {
		return "", datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return orderId, nil
}
