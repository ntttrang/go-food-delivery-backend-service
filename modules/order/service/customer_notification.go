package service

import (
	"context"
	"fmt"
	"log"

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Order Status Update"
		pushBody := fmt.Sprintf("Your order %s is now %s", orderID, s.getStateDisplayName(newState))
		data := map[string]interface{}{
			"orderID":  orderID,
			"newState": newState,
			"type":     "order_state_change",
		}
		if err := s.sendPushNotification(ctx, userID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to customer: %v", err)
			// Don't return error for push notification failure
		}
	}

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Shipper Assigned"
		pushBody := fmt.Sprintf("A shipper has been assigned to your order %s", orderID)
		data := map[string]interface{}{
			"orderID":   orderID,
			"shipperID": shipperID,
			"type":      "shipper_assignment",
		}
		if err := s.sendPushNotification(ctx, userID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to customer: %v", err)
			// Don't return error for push notification failure
		}
	}

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Payment Status Update"
		var pushBody string
		if paymentStatus == "paid" {
			pushBody = fmt.Sprintf("Payment confirmed for order %s. Your order is being prepared!", orderID)
		} else {
			pushBody = fmt.Sprintf("Payment status for order %s: %s", orderID, paymentStatus)
		}
		data := map[string]interface{}{
			"orderID":       orderID,
			"paymentStatus": paymentStatus,
			"type":          "payment_status_change",
		}
		if err := s.sendPushNotification(ctx, userID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to customer: %v", err)
			// Don't return error for push notification failure
		}
	}

	return nil
}

// notifyCustomerOrderCreated sends confirmation to customer about order creation
func (s *OrderNotificationService) notifyCustomerOrderCreated(ctx context.Context, userID, orderID, restaurantID string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error {
	// Get customer email
	userIdUUID, _ := uuid.Parse(userID)
	_, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userIdUUID})
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	email := "minhtrang.2106@gmail.com" //userMap[userIdUUID].Email

	// Create notification message
	subject, body := s.createCustomerOrderCreatedMessage(orderID, restaurantID, order, tracking, orderDetails)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to customer: %w", err)
	}

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Order Confirmed!"
		pushBody := fmt.Sprintf("Your order %s has been placed successfully. Total: $%.2f", orderID, order.TotalPrice)
		data := map[string]interface{}{
			"orderID":      orderID,
			"restaurantID": restaurantID,
			"totalPrice":   order.TotalPrice,
			"type":         "order_created",
		}
		if err := s.sendPushNotification(ctx, userID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to customer: %v", err)
			// Don't return error for push notification failure
		}
	}

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

	// Send push notification if available
	if s.pushSvc != nil {
		title := "Order Cancelled"
		pushBody := fmt.Sprintf("Your order %s has been cancelled. Reason: %s", orderID, reason)
		data := map[string]interface{}{
			"orderID":    orderID,
			"reason":     reason,
			"totalPrice": order.TotalPrice,
			"type":       "order_cancelled",
		}
		if err := s.sendPushNotification(ctx, userID, title, pushBody, data); err != nil {
			log.Printf("Failed to send push notification to customer: %v", err)
			// Don't return error for push notification failure
		}
	}

	return nil
}

// Customer message creation methods

// createCustomerStateChangeMessage creates email subject and body for customer notifications
func (s *OrderNotificationService) createCustomerStateChangeMessage(orderID, oldState, newState string) (string, string) {
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
