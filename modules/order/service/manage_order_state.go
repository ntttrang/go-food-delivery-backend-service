package service

import (
	"context"
	"time"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
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
	PaymentStatusPending = "PENDING"
	PaymentStatusPaid    = "PAID"
)

// StateTransitionRequest represents a request to change order state
type StateTransitionRequest struct {
	OrderID       string  `json:"orderId"`
	NewState      string  `json:"newState"`
	ShipperID     *string `json:"shipperId,omitempty"`
	PaymentStatus *string `json:"paymentStatus,omitempty"`
	UpdatedBy     string  `json:"updatedBy"` // User ID who is making the update
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
}

// Service
type OrderStateManagementService struct {
	repo                IOrderStateRepo
	notificationService IOrderNotificationService
}

func NewOrderStateManagementService(
	repo IOrderStateRepo,
	notificationService IOrderNotificationService,
) *OrderStateManagementService {
	return &OrderStateManagementService{
		repo:                repo,
		notificationService: notificationService,
	}
}

// ValidateStateTransition validates if the state transition is allowed
func (s *OrderStateManagementService) ValidateStateTransition(currentState, newState string) error {
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

	for _, allowedState := range allowedStates {
		if allowedState == newState {
			return nil
		}
	}

	return datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidOrderState).WithDebug("invalid state transition from " + currentState + " to " + newState)
}

// TransitionOrderState handles order state transitions
func (s *OrderStateManagementService) TransitionOrderState(ctx context.Context, req *StateTransitionRequest) error {
	// Get current order
	order, tracking, _, err := s.repo.FindById(ctx, req.OrderID)
	if err != nil {
		if err == ordermodel.ErrOrderNotFound {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Validate state transition
	if err := s.ValidateStateTransition(tracking.State, req.NewState); err != nil {
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
		if tracking.PaymentMethod == "cash" && tracking.PaymentStatus == PaymentStatusPending {
			tracking.PaymentStatus = PaymentStatusPaid
		}

	case StateCancelled:
		// Handle cancellation logic
		// For paid orders, might need to process refunds
		break
	}

	// Update payment status if provided
	if req.PaymentStatus != nil {
		if *req.PaymentStatus != PaymentStatusPending && *req.PaymentStatus != PaymentStatusPaid {
			return datatype.ErrBadRequest.WithError("invalid payment status")
		}
		tracking.PaymentStatus = *req.PaymentStatus
	}

	// Save changes
	if err := s.repo.Update(ctx, order, tracking); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Send notifications
	if s.notificationService != nil {
		// Notify state change
		if err := s.notificationService.NotifyOrderStateChange(ctx, req.OrderID, oldState, req.NewState); err != nil {
			// Log error but don't fail the operation
			// In production, you might want to use a message queue for reliability
		}

		// Notify shipper assignment
		if req.ShipperID != nil && order.ShipperID != nil {
			if err := s.notificationService.NotifyShipperAssignment(ctx, req.OrderID, *order.ShipperID); err != nil {
				// Log error but don't fail the operation
			}
		}

		// Notify payment status change
		if req.PaymentStatus != nil {
			if err := s.notificationService.NotifyPaymentStatusChange(ctx, req.OrderID, *req.PaymentStatus); err != nil {
				// Log error but don't fail the operation
			}
		}
	}

	return nil
}

// AssignShipper assigns a shipper to an order
func (s *OrderStateManagementService) AssignShipper(ctx context.Context, orderID string, shipperID string) error {
	// Get current order
	order, tracking, _, err := s.repo.FindById(ctx, orderID)
	if err != nil {
		if err == ordermodel.ErrOrderNotFound {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Check if order is in a state where shipper can be assigned
	if tracking.State != StateWaitingForShipper && tracking.State != StatePreparing {
		return datatype.ErrBadRequest.WithError("shipper can only be assigned to orders waiting for shipper or preparing")
	}

	// Assign shipper
	order.ShipperID = &shipperID
	order.UpdatedAt = time.Now()

	// If order is waiting for shipper, move to preparing
	if tracking.State == StateWaitingForShipper {
		tracking.State = StatePreparing
		tracking.UpdatedAt = time.Now()
	}

	// Save changes
	if err := s.repo.Update(ctx, order, tracking); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Send notification
	if s.notificationService != nil {
		if err := s.notificationService.NotifyShipperAssignment(ctx, orderID, shipperID); err != nil {
			// Log error but don't fail the operation
		}
	}

	return nil
}

// UpdatePaymentStatus updates the payment status of an order
func (s *OrderStateManagementService) UpdatePaymentStatus(ctx context.Context, orderID string, paymentStatus string) error {
	// Validate payment status
	if paymentStatus != PaymentStatusPending && paymentStatus != PaymentStatusPaid {
		return datatype.ErrBadRequest.WithError("invalid payment status")
	}

	// Get current order
	_, tracking, _, err := s.repo.FindById(ctx, orderID)
	if err != nil {
		if err == ordermodel.ErrOrderNotFound {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Update payment status
	tracking.PaymentStatus = paymentStatus
	tracking.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, nil, tracking); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Send notification
	if s.notificationService != nil {
		if err := s.notificationService.NotifyPaymentStatusChange(ctx, orderID, paymentStatus); err != nil {
			// Log error but don't fail the operation
		}
	}

	return nil
}
