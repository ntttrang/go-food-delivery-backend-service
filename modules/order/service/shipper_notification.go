package service

import (
	"context"
	"fmt"
	"log"

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "New Order Assignment"
		pushBody := fmt.Sprintf("You have been assigned to order %s", orderID)
		data := map[string]interface{}{
			"orderID":      orderID,
			"restaurantID": restaurantID,
			"type":         "shipper_assignment",
		}
		if err := s.sendPushNotification(ctx, shipperID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to shipper: %v", err)
			// Don't return error for push notification failure
		}
	}

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Payment Confirmed"
		pushBody := fmt.Sprintf("Payment confirmed for order %s. Order is ready for processing!", orderID)
		data := map[string]interface{}{
			"orderID":       orderID,
			"paymentStatus": paymentStatus,
			"type":          "payment_status_change",
		}
		if err := s.sendPushNotification(ctx, shipperID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to shipper: %v", err)
			// Don't return error for push notification failure
		}
	}

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "New Order Available"
		pushBody := fmt.Sprintf("New order %s assigned to you. Value: $%.2f", orderID, order.TotalPrice)
		data := map[string]interface{}{
			"orderID":      orderID,
			"restaurantID": restaurantID,
			"totalPrice":   order.TotalPrice,
			"type":         "order_created",
		}
		if err := s.sendPushNotification(ctx, shipperID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to shipper: %v", err)
			// Don't return error for push notification failure
		}
	}

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Order Cancelled"
		pushBody := fmt.Sprintf("Order %s has been cancelled. You are no longer assigned to this order.", orderID)
		data := map[string]interface{}{
			"orderID":    orderID,
			"reason":     reason,
			"totalPrice": order.TotalPrice,
			"type":       "order_cancelled",
		}
		if err := s.sendPushNotification(ctx, shipperID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to shipper: %v", err)
			// Don't return error for push notification failure
		}
	}

	return nil
}
