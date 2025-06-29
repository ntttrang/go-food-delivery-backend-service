package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared"
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

type IEvtPublisher interface {
	Publish(ctx context.Context, topic string, evt *datatype.AppEvent) error
}

// CreateFromCartCommandHandler handles creating orders from cart
type CreateFromCartCommandHandler struct {
	createHandler         *CreateCommandHandler
	cartConversionService ICartConversionService
	paymentService        *PaymentProcessingService
	inventoryService      *InventoryCheckingService
	notificationService   *OrderNotificationService
	evtPublisher          IEvtPublisher
}

func NewCreateFromCartCommandHandler(
	createHandler *CreateCommandHandler,
	cartConversionService ICartConversionService,
	paymentService *PaymentProcessingService,
	inventoryService *InventoryCheckingService,
	notificationService *OrderNotificationService,
	evtPublisher IEvtPublisher,
) *CreateFromCartCommandHandler {
	return &CreateFromCartCommandHandler{
		createHandler:         createHandler,
		cartConversionService: cartConversionService,
		paymentService:        paymentService,
		inventoryService:      inventoryService,
		notificationService:   notificationService,
		evtPublisher:          evtPublisher,
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

	// Create the order using the standard flow
	orderId, err := s.createHandler.Execute(ctx, orderData)
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
			log.Print("Payment failed! \n")
			return "", datatype.ErrBadRequest.WithWrap(ordermodel.ErrPaymentFailed).WithDebug(paymentResult.ErrorMessage)
		}
	}

	// Mark cart as processed
	if err := s.cartConversionService.MarkCartAsProcessed(ctx, cartID); err != nil {
		log.Print("update cart status = PROCESSED after order created \n")
	}

	// Send notifications
	orderCreatedMsg := map[string]interface{}{
		"orderId":      orderId,
		"userId":       data.UserID,
		"restaurantId": orderData.RestaurantID,
	}
	go func() {
		log.Println("Publish msg: ORDER CREATED")
		defer shared.Recover()

		evt := datatype.NewAppEvent(
			datatype.WithTopic(datatype.EvtNotifyOrderCreate),
			datatype.WithData(orderCreatedMsg),
		)
		if err := s.evtPublisher.Publish(ctx, evt.Topic, evt); err != nil {
			log.Println("Failed to publish event", err)
		}
	}()

	return orderId, nil
}
