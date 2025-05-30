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
	CardID          *string                `json:"cardId,omitempty"` // Required for card payments
	OrderDetails    []OrderDetailCreateDto `json:"orderDetails"`
}

// New DTO for creating order from cart
type OrderCreateFromCartDto struct {
	UserID          string              `json:"-"` // Get from token
	CartID          string              `json:"cartId"`
	DeliveryAddress *ordermodel.Address `json:"deliveryAddress"`
	PaymentMethod   string              `json:"paymentMethod"`
	CardID          *string             `json:"cardId,omitempty"` // Required for card payments
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

	if o.PaymentMethod == "" {
		return ordermodel.ErrPaymentMethodRequired
	}

	// Validate payment method and card requirement
	if o.PaymentMethod == "card" && o.CardID == nil {
		return ordermodel.ErrCardIdRequired
	}

	if o.DeliveryAddress == nil {
		return ordermodel.ErrDeliveryAddressRequired
	}

	return nil
}

func (o *OrderCreateFromCartDto) Validate() error {
	if o.UserID == "" {
		return ordermodel.ErrUserIdRequired
	}

	if o.CartID == "" {
		return ordermodel.ErrCartIdRequired
	}

	if o.PaymentMethod == "" {
		return ordermodel.ErrPaymentMethodRequired
	}

	// Validate payment method and card requirement
	if o.PaymentMethod == "card" && o.CardID == nil {
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
	repo                  ICreateOrderRepository
	cartConversionService *CartToOrderConversionService
	paymentService        *PaymentProcessingService
	inventoryService      *InventoryCheckingService
	notificationService   *OrderNotificationService
}

func NewCreateCommandHandler(
	repo ICreateOrderRepository,
	cartConversionService *CartToOrderConversionService,
	paymentService *PaymentProcessingService,
	inventoryService *InventoryCheckingService,
	notificationService *OrderNotificationService,
) *CreateCommandHandler {
	return &CreateCommandHandler{
		repo:                  repo,
		cartConversionService: cartConversionService,
		paymentService:        paymentService,
		inventoryService:      inventoryService,
		notificationService:   notificationService,
	}
}

// Enhanced CreateCommandHandler for backward compatibility
func NewCreateCommandHandlerSimple(repo ICreateOrderRepository) *CreateCommandHandler {
	return &CreateCommandHandler{
		repo: repo,
	}
}

// Execute creates an order with manual order details
// NOTE: This method is restricted to ADMIN users only via API layer
// Use ExecuteFromCart for standard customer order creation
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

// ExecuteFromCart creates an order from a cart
func (s *CreateCommandHandler) ExecuteFromCart(ctx context.Context, data *OrderCreateFromCartDto) (string, error) {
	if err := data.Validate(); err != nil {
		return "", datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Check if all required services are available
	if s.cartConversionService == nil {
		return "", datatype.ErrInternalServerError.WithError("cart conversion service not available")
	}

	userID, err := uuid.Parse(data.UserID)
	if err != nil {
		return "", datatype.ErrBadRequest.WithError("invalid user ID format")
	}

	cartID, err := uuid.Parse(data.CartID)
	if err != nil {
		return "", datatype.ErrBadRequest.WithError("invalid cart ID format")
	}

	// Validate cart can be converted to order
	if err := s.cartConversionService.ValidateCartForOrder(ctx, cartID, userID); err != nil {
		return "", err
	}

	// Convert cart to order data
	orderData, err := s.cartConversionService.ConvertCartToOrderData(ctx, cartID, userID)
	if err != nil {
		return "", err
	}

	// Set delivery address and payment method from request
	orderData.DeliveryAddress = data.DeliveryAddress
	orderData.PaymentMethod = data.PaymentMethod
	orderData.CardID = data.CardID

	// Check inventory if service is available
	if s.inventoryService != nil {
		restaurantID, err := uuid.Parse(orderData.RestaurantID)
		if err != nil {
			return "", datatype.ErrBadRequest.WithError("invalid restaurant ID format")
		}

		// Convert order details to inventory items
		var inventoryItems []OrderItem
		for _, detail := range orderData.OrderDetails {
			foodID, err := uuid.Parse(detail.FoodOrigin.Id)
			if err != nil {
				return "", datatype.ErrBadRequest.WithError("invalid food ID format")
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
			Amount:        orderData.TotalPrice,
			PaymentMethod: data.PaymentMethod,
			CardID:        data.CardID,
		}

		// Validate payment method first
		if err := s.paymentService.ValidatePaymentMethod(ctx, paymentReq); err != nil {
			return "", err
		}
	}

	// Create the order using the standard flow
	orderId, err := s.Execute(ctx, orderData)
	if err != nil {
		return "", err
	}

	// Process payment after order creation
	if s.paymentService != nil {
		paymentReq := &PaymentRequest{
			OrderID:       orderId,
			UserID:        data.UserID,
			Amount:        orderData.TotalPrice,
			PaymentMethod: data.PaymentMethod,
			CardID:        data.CardID,
		}

		paymentResult, err := s.paymentService.ProcessPayment(ctx, paymentReq)
		if err != nil {
			// TODO: In a real implementation, you might want to cancel the order here
			// or mark it as payment failed
			return "", datatype.ErrInternalServerError.WithWrap(err).WithDebug("payment processing failed")
		}

		if !paymentResult.Success {
			// TODO: Handle payment failure
			return "", datatype.ErrBadRequest.WithWrap(ordermodel.ErrPaymentFailed).WithDebug(paymentResult.ErrorMessage)
		}
	}

	// Mark cart as processed
	if err := s.cartConversionService.MarkCartAsProcessed(ctx, cartID); err != nil {
		// Log error but don't fail the order creation
		// The order was successfully created, cart processing is secondary
	}

	// Send notifications
	if s.notificationService != nil {
		if err := s.notificationService.NotifyOrderCreated(ctx, orderId, data.UserID, orderData.RestaurantID); err != nil {
			// Log error but don't fail the operation
		}
	}

	return orderId, nil
}
