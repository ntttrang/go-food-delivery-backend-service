package service

import (
	"context"
	"log"
	"time"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// OrderState constants
const (
	StateWaitingForShipper = "waiting_for_shipper"
	StatePreparing         = "preparing"
	StateOnTheWay          = "on_the_way"
	StateDelivered         = "delivered"
	StateCancelled         = "cancel"
)

// PaymentStatus constants
const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
)

// StateTransitionRequest represents a request to change order state
type StateTransitionRequest struct {
	OrderID            string  `json:"orderId"`
	NewState           string  `json:"newState"`
	ShipperID          *string `json:"shipperId,omitempty"`
	PaymentStatus      *string `json:"paymentStatus,omitempty"`
	UpdatedBy          string  `json:"-"`                            // User ID who is making the update - Get from Requester context
	CancellationReason *string `json:"cancellationReason,omitempty"` // Required when cancelling
}

// Repository interface
type IOrderStateRepo interface {
	FindById(ctx context.Context, id string) (*ordermodel.Order, *ordermodel.OrderTracking, []ordermodel.OrderDetail, error)
	Update(ctx context.Context, order *ordermodel.Order, tracking *ordermodel.OrderTracking) error
}

// Notification interface for future implementation
type IOrderNotificationService interface {
	NotifyOrderStateChange(ctx context.Context, orderID string, oldState string, newState string) error
	NotifyShipperAssignment(ctx context.Context, orderID string, shipperID string) error
	NotifyPaymentStatusChange(ctx context.Context, orderID string, paymentStatus string) error
	NotifyOrderCancelled(ctx context.Context, orderID string, reason string) error
}

// Service
type OrderStateManagementService struct {
	repo                IOrderStateRepo
	notificationService IOrderNotificationService
	evtPublisher        IEvtPublisher
}

func NewOrderStateManagementService(
	repo IOrderStateRepo,
	notificationService IOrderNotificationService,
	evtPublisher IEvtPublisher,
) *OrderStateManagementService {
	return &OrderStateManagementService{
		repo:                repo,
		notificationService: notificationService,
		// refundService and inventoryService are optional for now
		evtPublisher: evtPublisher,
	}
}

// validateStateTransition validates if the state transition is allowed
func (s *OrderStateManagementService) validateStateTransition(currentState, newState string) error {
	validTransitions := map[string][]string{
		StateWaitingForShipper: {StatePreparing, StateCancelled},
		StatePreparing:         {StateOnTheWay, StateCancelled},
		StateOnTheWay:          {StateDelivered, StateCancelled},
		StateDelivered:         {}, // Terminal state
		StateCancelled:         {}, // Terminal state
	}

	allowedStates, exists := validTransitions[currentState]
	if !exists {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidOrderState).WithDebug("invalid current state: " + currentState)
	}

	// Check if newState is in allowedStates
	for _, allowedState := range allowedStates {
		if allowedState == newState {
			return nil
		}
	}

	return datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidOrderState).WithDebug("invalid state transition from " + currentState + " to " + newState)
}

// Execute handles order state transitions
func (s *OrderStateManagementService) Execute(ctx context.Context, req *StateTransitionRequest) error {
	// Get current order
	order, tracking, _, err := s.repo.FindById(ctx, req.OrderID)
	if err != nil {
		if err == ordermodel.ErrOrderNotFound {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Validate state transition
	if err := s.validateStateTransition(tracking.State, req.NewState); err != nil {
		return err
	}

	// Store old state for notification
	oldState := tracking.State

	// Update tracking state
	tracking.State = req.NewState
	tracking.UpdatedAt = time.Now()

	// Handle specific state transitions
	switch req.NewState {
	case StatePreparing:
		// When order moves to preparing, it should have a shipper assigned
		if req.ShipperID != nil {
			order.ShipperID = req.ShipperID
			order.UpdatedAt = time.Now()
		}

	case StateOnTheWay:
		// Ensure shipper is assigned
		if order.ShipperID == nil {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrShipperRequired).WithDebug("shipper must be assigned before order can be on the way")
		}
		// Set estimated delivery time (example: 30 minutes from now)
		tracking.EstimatedTime = 30

	case StateDelivered:
		// Set actual delivery time
		tracking.DeliveryTime = int(time.Since(tracking.CreatedAt).Minutes())
		// For cash payments, mark as paid when delivered
		if tracking.PaymentMethod == MethodCash && tracking.PaymentStatus == PaymentStatusPending {
			tracking.PaymentStatus = PaymentStatusPaid
		}

	case StateCancelled:
		// Validate cancellation reason is provided
		if req.CancellationReason == nil || *req.CancellationReason == "" {
			return datatype.ErrBadRequest.WithError("cancellation reason is required when cancelling an order")
		}

		// Handle refund for paid orders
		if tracking.PaymentStatus == PaymentStatusPaid {
			// Step 1: Process refund based on payment method ( TBD)
			// Step 2: Update payment status to indicate refund is being processed
			tracking.PaymentStatus = PaymentStatusPending
		}

		// Step 3: Restore inventory if inventory service is available (TBD)

		// Clear shipper assignment if order was assigned
		if order.ShipperID != nil {
			order.ShipperID = nil
			order.UpdatedAt = time.Now()
		}

		// Set cancellation timestamp (using DeliveryTime field to track cancellation time)
		tracking.CancelReason = *req.CancellationReason
		tracking.DeliveryTime = int(time.Since(tracking.CreatedAt).Minutes())
	}

	// Update payment status if provided
	if req.PaymentStatus != nil {
		if *req.PaymentStatus != PaymentStatusPending && *req.PaymentStatus != PaymentStatusPaid {
			return datatype.ErrBadRequest.WithError("invalid payment status")
		}
		tracking.PaymentStatus = *req.PaymentStatus
	}
	order.UpdatedBy = &req.UpdatedBy
	tracking.UpdatedBy = &req.UpdatedBy

	// Save changes
	if err := s.repo.Update(ctx, order, tracking); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Send notifications
	// Change state
	orderStateChangeMsg := map[string]interface{}{
		"orderId":  req.OrderID,
		"oldState": oldState,
		"newState": req.NewState,
	}
	go func() {
		log.Println("Publish msg: ORDER STATE CHANGE")
		defer shared.Recover()

		evt := datatype.NewAppEvent(
			datatype.WithTopic(datatype.EvtNotifyOrderStateChange),
			datatype.WithData(orderStateChangeMsg),
		)

		if err := s.evtPublisher.Publish(ctx, evt.Topic, evt); err != nil {
			log.Println("Failed to publish event", err)
		}
	}()

	// Notify cancellation with reason
	if req.NewState == StateCancelled && req.CancellationReason != nil {
		orderCancelMsg := map[string]interface{}{
			"orderId":      req.OrderID,
			"cancelReason": *req.CancellationReason,
		}
		go func() {
			log.Println("Publish msg: CANCEL ORDER")
			defer shared.Recover()

			evt := datatype.NewAppEvent(
				datatype.WithTopic(datatype.EvtNotifyOrderCancel),
				datatype.WithData(orderCancelMsg),
			)

			if err := s.evtPublisher.Publish(ctx, evt.Topic, evt); err != nil {
				log.Println("Failed to publish event", err)
			}
		}()
	}

	// Notify shipper assignment
	if req.ShipperID != nil && order.ShipperID != nil {
		assignShipperMsg := map[string]interface{}{
			"orderId":   req.OrderID,
			"shipperId": *order.ShipperID,
		}
		go func() {
			log.Println("Publish msg: ASSIGN SHIPPER")
			defer shared.Recover()

			evt := datatype.NewAppEvent(
				datatype.WithTopic(datatype.EvtNotifyShipperAssign),
				datatype.WithData(assignShipperMsg),
			)

			if err := s.evtPublisher.Publish(ctx, evt.Topic, evt); err != nil {
				log.Println("Failed to publish event", err)
			}
		}()
	}

	// Notify payment status change
	if req.PaymentStatus != nil {
		paymentStatusChangeMsg := map[string]interface{}{
			"orderId":       req.OrderID,
			"paymentStatus": req.PaymentStatus,
		}
		go func() {
			log.Println("Publish msg: CHANGE PAYMENT STATUS")
			defer shared.Recover()

			evt := datatype.NewAppEvent(
				datatype.WithTopic(datatype.EvtNotifyPaymentStatusChange),
				datatype.WithData(paymentStatusChangeMsg),
			)

			if err := s.evtPublisher.Publish(ctx, evt.Topic, evt); err != nil {
				log.Println("Failed to publish event", err)
			}
		}()

	}

	return nil
}
