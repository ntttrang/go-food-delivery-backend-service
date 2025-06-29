package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// Shipper notification methods

// notifyShipperStateChange sends notification to shipper about order state change
func (s *OrderNotificationService) notifyShipperStateChange(ctx context.Context, shipperID, orderID, oldState, newState string) error {
	// Get shipper email
	shipperIdUUID, _ := uuid.Parse(shipperID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{shipperIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get shipper email: %w", err)
	}

	email := userMap[shipperIdUUID].Email
	// Create notification message
	subject, body := s.createShipperStateChangeMessage(orderID, oldState, newState)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to shipper: %w", err)
	}

	return nil
}

// notifyShipperAssignment sends notification to shipper about order assignment
func (s *OrderNotificationService) notifyShipperAssignment(ctx context.Context, shipperID, orderID, restaurantID string) error {
	// Get shipper email
	shipperIdUUID, _ := uuid.Parse(shipperID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{shipperIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get shipper: %w", err)
	}

	email := userMap[shipperIdUUID].Email

	// Create notification message
	subject, body := s.createShipperAssignmentMessage(orderID, restaurantID)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to shipper: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyShipperPaymentStatusChange sends notification to shipper about payment status change
func (s *OrderNotificationService) notifyShipperPaymentStatusChange(ctx context.Context, shipperID, orderID, paymentStatus string) error {
	// Get shipper email
	shipperIdUUID, _ := uuid.Parse(shipperID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{shipperIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get shipper: %w", err)
	}

	email := userMap[shipperIdUUID].Email

	// Create notification message
	subject, body := s.createShipperPaymentStatusChangeMessage(orderID, paymentStatus)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to shipper: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyShipperOrderCreated sends notification to shipper about new order (if already assigned)
func (s *OrderNotificationService) notifyShipperOrderCreated(ctx context.Context, shipperID, orderID, restaurantID string, order *ordermodel.Order, tracking *ordermodel.OrderTracking) error {
	// Get shipper email
	shipperIdUUID, _ := uuid.Parse(shipperID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{shipperIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get shipper: %w", err)
	}

	email := userMap[shipperIdUUID].Email

	// Create notification message
	subject, body := s.createShipperOrderCreatedMessage(orderID, restaurantID, order, tracking)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to shipper: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyShipperOrderCancelled sends cancellation notification to shipper
func (s *OrderNotificationService) notifyShipperOrderCancelled(ctx context.Context, shipperID, orderID, reason string, order *ordermodel.Order, tracking *ordermodel.OrderTracking) error {
	// Get shipper email
	shipperIdUUID, _ := uuid.Parse(shipperID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{shipperIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get shipper: %w", err)
	}

	email := userMap[shipperIdUUID].Email

	// Create notification message
	subject, body := s.createShipperOrderCancelledMessage(orderID, reason, order, tracking)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to shipper: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}
