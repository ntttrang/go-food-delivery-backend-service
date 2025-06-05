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
	CardID          string                 `json:"cardId,omitempty"` // Required for card payments
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

// Payment method constants
const (
	MethodCash       = "cash"
	MethodCreditCard = "credit_card"
	MethodDebitCard  = "debit_card"
)

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

	if o.PaymentMethod == "" {
		return ordermodel.ErrPaymentMethodRequired
	}

	// Validate payment method and card requirement
	if (o.PaymentMethod == MethodCreditCard || o.PaymentMethod == MethodDebitCard) && o.CardID == "" {
		return ordermodel.ErrCardIdRequired
	}

	if o.DeliveryAddress == nil {
		return ordermodel.ErrDeliveryAddressRequired
	}

	return nil
}

// Initialize service
type ICreateOrderRepository interface {
	Insert(ctx context.Context, order *ordermodel.Order, orderTracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error
}

type CreateCommandHandler struct {
	repo                ICreateOrderRepository
	paymentService      *PaymentProcessingService
	inventoryService    *InventoryCheckingService
	notificationService *OrderNotificationService
}

func NewCreateCommandHandler(
	repo ICreateOrderRepository,
	paymentService *PaymentProcessingService,
	inventoryService *InventoryCheckingService,
	notificationService *OrderNotificationService,
) *CreateCommandHandler {
	return &CreateCommandHandler{
		repo:                repo,
		paymentService:      paymentService,
		inventoryService:    inventoryService,
		notificationService: notificationService,
	}
}

func (s *CreateCommandHandler) Execute(ctx context.Context, data *OrderCreateDto) (string, error) {
	// Validate request
	if err := data.Validate(); err != nil {
		return "", datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Check inventory if service is available
	if s.inventoryService != nil {
		restaurantID, err := uuid.Parse(data.RestaurantID)
		if err != nil {
			return "", datatype.ErrBadRequest.WithError(ordermodel.ErrInvalidRestaurantIdFormat.Error())
		}

		// Convert order details to inventory items
		var inventoryItems []OrderItem
		for _, detail := range data.OrderDetails {
			foodID, err := uuid.Parse(detail.FoodOrigin.Id)
			if err != nil {
				return "", datatype.ErrBadRequest.WithError(ordermodel.ErrInvalidFoodIdFormat.Error())
			}
			inventoryItems = append(inventoryItems, OrderItem{
				FoodID:   foodID,
				Quantity: detail.Quantity,
			})
		}

		// Check inventory
		if err := s.inventoryService.CheckOrderInventory(ctx, restaurantID, inventoryItems); err != nil {
			return "", err
		}
	}

	// Process payment if service is available
	if s.paymentService != nil {
		paymentReq := &PaymentRequest{
			UserID:        data.UserID,
			Amount:        data.TotalPrice,
			PaymentMethod: data.PaymentMethod,
			CardID:        &data.CardID,
		}

		// Validate payment method first
		if err := s.paymentService.ValidatePaymentMethod(ctx, paymentReq); err != nil {
			return "", err
		}
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
		State:           StateWaitingForShipper,
		PaymentStatus:   PaymentStatusPending,
		PaymentMethod:   data.PaymentMethod,
		CardId:          &data.CardID,
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
