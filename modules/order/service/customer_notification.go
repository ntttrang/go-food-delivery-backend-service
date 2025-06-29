package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// Customer notification methods

// notifyCustomerStateChange sends notification to customer about order state change
func (s *OrderNotificationService) notifyCustomerStateChange(ctx context.Context, userID, orderID, oldState, newState string) error {
	// Get customer email
	userIdUUID, _ := uuid.Parse(userID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	email := userMap[userIdUUID].Email
	// Create notification message
	subject, body := s.createCustomerStateChangeMessage(orderID, oldState, newState)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to customer: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyCustomerShipperAssignment sends notification to customer about shipper assignment
func (s *OrderNotificationService) notifyCustomerShipperAssignment(ctx context.Context, userID, orderID, shipperID string) error {
	// Get customer email
	userIdUUID, _ := uuid.Parse(userID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	email := userMap[userIdUUID].Email

	// Create notification message
	subject, body := s.createCustomerShipperAssignmentMessage(orderID, shipperID)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to customer: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyCustomerPaymentStatusChange sends notification to customer about payment status change
func (s *OrderNotificationService) notifyCustomerPaymentStatusChange(ctx context.Context, userID, orderID, paymentStatus string) error {
	// Get customer email
	userIdUUID, _ := uuid.Parse(userID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	email := userMap[userIdUUID].Email

	// Create notification message
	subject, body := s.createCustomerPaymentStatusChangeMessage(orderID, paymentStatus)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to customer: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyCustomerOrderCreated sends confirmation to customer about order creation
func (s *OrderNotificationService) notifyCustomerOrderCreated(ctx context.Context, userID, orderID, restaurantID string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error {
	// Get customer email
	userIdUUID, _ := uuid.Parse(userID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	email := userMap[userIdUUID].Email

	// Create notification message
	subject, body := s.createCustomerOrderCreatedMessage(orderID, restaurantID, order, tracking, orderDetails)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to customer: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// notifyCustomerOrderCancelled sends cancellation notification to customer
func (s *OrderNotificationService) notifyCustomerOrderCancelled(ctx context.Context, userID, orderID, reason string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error {
	// Get customer email
	userIdUUID, _ := uuid.Parse(userID)
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	email := userMap[userIdUUID].Email

	// Create notification message
	subject, body := s.createCustomerOrderCancelledMessage(orderID, reason, order, tracking, orderDetails)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to customer: %w", err)
	}

	// TODO: Send push notification if available ( TBD)

	return nil
}

// Customer message creation methods

// createCustomerStateChangeMessage creates email subject and body for customer notifications
func (s *OrderNotificationService) createCustomerStateChangeMessage(orderID, _ /* oldState */, newState string) (string, string) {
	subject := fmt.Sprintf("Order %s: %s", orderID, s.getStateDisplayName(newState))

	var body string
	switch newState {
	case "preparing":
		body = fmt.Sprintf(`Order %s is being prepared 🍳

We'll notify you when ready for delivery.

Food Delivery Team`, orderID)

	case "on_the_way":
		body = fmt.Sprintf(`Order %s is on the way! 🚗

Please be available to receive your order.

Food Delivery Team`, orderID)

	case "delivered":
		body = fmt.Sprintf(`Order %s delivered! ✅

Enjoy your meal!

Food Delivery Team`, orderID)

	case "cancel":
		body = fmt.Sprintf(`Order %s cancelled ❌

Contact support for questions.

Food Delivery Team`, orderID)

	default:
		body = fmt.Sprintf(`Order %s: %s

Food Delivery Team`, orderID, s.getStateDisplayName(newState))
	}

	return subject, body
}
