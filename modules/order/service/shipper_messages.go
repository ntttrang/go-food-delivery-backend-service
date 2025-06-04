package service

import (
	"fmt"
	"strings"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// Shipper message creation methods

// createShipperStateChangeMessage creates email subject and body for shipper notifications
func (s *OrderNotificationService) createShipperStateChangeMessage(orderID, oldState, newState string) (string, string) {
	subject := fmt.Sprintf("Order %s: %s", orderID, s.getStateDisplayName(newState))

	var body string
	switch newState {
	case "preparing":
		body = fmt.Sprintf(`Order %s: Being prepared

Be ready for pickup notification.

Food Delivery Platform`, orderID)

	case "on_the_way":
		body = fmt.Sprintf(`Order %s: On the way 🚗

Ensure safe delivery.

Food Delivery Platform`, orderID)

	case "delivered":
		body = fmt.Sprintf(`Order %s: Delivered ✅

Thank you for your service!

Food Delivery Platform`, orderID)

	case "cancel":
		body = fmt.Sprintf(`Order %s: Cancelled

You are no longer assigned.

Food Delivery Platform`, orderID)

	default:
		body = fmt.Sprintf(`Order %s: %s

Food Delivery Platform`, orderID, s.getStateDisplayName(newState))
	}

	return subject, body
}

// createShipperAssignmentMessage creates email subject and body for shipper assignment notifications
func (s *OrderNotificationService) createShipperAssignmentMessage(orderID, restaurantID string) (string, string) {
	subject := fmt.Sprintf("New Assignment - Order %s", orderID)

	body := fmt.Sprintf(`Order %s assigned to you

Prepare for pickup when ready.
Check app for details.

Food Delivery Platform`, orderID)

	return subject, body
}

// createShipperPaymentStatusChangeMessage creates email subject and body for shipper payment status notifications
func (s *OrderNotificationService) createShipperPaymentStatusChangeMessage(orderID, paymentStatus string) (string, string) {
	subject := fmt.Sprintf("Payment %s - Order %s", strings.ToUpper(paymentStatus), orderID)

	var body string
	switch paymentStatus {
	case "paid":
		body = fmt.Sprintf(`Order %s: Payment confirmed ✓

Order ready for processing.

Food Delivery Platform`, orderID)

	default:
		body = fmt.Sprintf(`Order %s: Payment %s

Wait for confirmation.

Food Delivery Platform`, orderID, strings.ToUpper(paymentStatus))
	}

	return subject, body
}

// createShipperOrderCreatedMessage creates email subject and body for shipper order creation notifications
func (s *OrderNotificationService) createShipperOrderCreatedMessage(orderID, restaurantID string, order *ordermodel.Order, tracking *ordermodel.OrderTracking) (string, string) {
	subject := fmt.Sprintf("New Order #%s", orderID)

	body := fmt.Sprintf(`Order #%s assigned

Value: $%.2f
Delivery: %d min
Address: %s

Wait for pickup notification.

Food Delivery Platform`, orderID, order.TotalPrice, tracking.EstimatedTime, string(tracking.DeliveryAddress))

	return subject, body
}

// createShipperOrderCancelledMessage creates email subject and body for shipper order cancellation notifications
func (s *OrderNotificationService) createShipperOrderCancelledMessage(orderID, reason string, order *ordermodel.Order, tracking *ordermodel.OrderTracking) (string, string) {
	subject := fmt.Sprintf("Order Cancelled #%s", orderID)

	body := fmt.Sprintf(`Order #%s cancelled

Reason: %s
Value: $%.2f

You are no longer assigned.
Check app for new orders.

Food Delivery Platform`, orderID, reason, order.TotalPrice)

	return subject, body
}
