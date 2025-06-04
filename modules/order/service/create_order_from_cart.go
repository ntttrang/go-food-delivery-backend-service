package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// DTO for creating order from cart
type OrderCreateFromCartDto struct {
	UserID          string              `json:"-"` // Get from token
	CartID          string              `json:"cartId"`
	DeliveryAddress *ordermodel.Address `json:"deliveryAddress"`
	PaymentMethod   string              `json:"paymentMethod"`
	CardID          *string             `json:"cardId,omitempty"` // Required for card payments
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
	if (o.PaymentMethod == MethodCreditCard || o.PaymentMethod == MethodDebitCard) && *o.CardID == "" {
		return ordermodel.ErrCardIdRequired
	}

	if o.DeliveryAddress == nil {
		return ordermodel.ErrDeliveryAddressRequired
	}

	return nil
}

// CreateFromCartCommandHandler handles creating orders from cart
type CreateFromCartCommandHandler struct {
	createHandler         *CreateCommandHandler
	cartConversionService ICartConversionService
	paymentService        *PaymentProcessingService
	inventoryService      *InventoryCheckingService
	notificationService   *OrderNotificationService
}

func NewCreateFromCartCommandHandler(
	createHandler *CreateCommandHandler,
	cartConversionService ICartConversionService,
	paymentService *PaymentProcessingService,
	inventoryService *InventoryCheckingService,
	notificationService *OrderNotificationService,
) *CreateFromCartCommandHandler {
	return &CreateFromCartCommandHandler{
		createHandler:         createHandler,
		cartConversionService: cartConversionService,
		paymentService:        paymentService,
		inventoryService:      inventoryService,
		notificationService:   notificationService,
	}
}

// ExecuteFromCart creates an order from a cart
func (s *CreateFromCartCommandHandler) ExecuteFromCart(ctx context.Context, data *OrderCreateFromCartDto) (string, error) {
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
	orderData.CardID = *data.CardID

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
	orderId, err := s.createHandler.Execute(ctx, orderData)
	if err != nil {
		return "", err
	}

	// TODO: Rework as event-driven
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
			log.Print("Payment failed! \n")
			return "", datatype.ErrBadRequest.WithWrap(ordermodel.ErrPaymentFailed).WithDebug(paymentResult.ErrorMessage)
		}
	}

	// Mark cart as processed
	if err := s.cartConversionService.MarkCartAsProcessed(ctx, cartID); err != nil {
		log.Print("update cart status = PROCESSED after order created \n")
	}

	// 	TODO: Rework as event-driven
	// Send notifications
	if s.notificationService != nil {
		if err := s.notificationService.NotifyOrderCreated(ctx, orderId, data.UserID, orderData.RestaurantID); err != nil {
			log.Print("Notify order created \n")
		}
	}

	return orderId, nil
}
