package service

import (
	"fmt"
	"strings"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// Restaurant message creation methods

// createRestaurantStateChangeMessage creates email subject and body for restaurant notifications
func (s *OrderNotificationService) createRestaurantStateChangeMessage(orderID, _ /* oldState */, newState string) (string, string) {
	subject := fmt.Sprintf("Order %s: %s", orderID, s.getStateDisplayName(newState))

	var body string
	switch newState {
	case "preparing":
		body = fmt.Sprintf(`Order %s: In preparation

Please prepare according to specifications.

Food Delivery Platform`, orderID)

	case "on_the_way":
		body = fmt.Sprintf(`Order %s: Picked up

Order is on the way to customer.

Food Delivery Platform`, orderID)

	case "delivered":
		body = fmt.Sprintf(`Order %s: Delivered ✅

Thank you for your service!

Food Delivery Platform`, orderID)

	case "cancel":
		body = fmt.Sprintf(`Order %s: Cancelled

Stop preparation if not completed.

Food Delivery Platform`, orderID)

	default:
		body = fmt.Sprintf(`Order %s: %s

Food Delivery Platform`, orderID, s.getStateDisplayName(newState))
	}

	return subject, body
}

// createRestaurantShipperAssignmentMessage creates email subject and body for restaurant shipper assignment notifications
func (s *OrderNotificationService) createRestaurantShipperAssignmentMessage(orderID, _ /* shipperID */ string) (string, string) {
	subject := fmt.Sprintf("Shipper Assigned - Order %s", orderID)

	body := fmt.Sprintf(`Order %s: Shipper assigned

Prepare order for pickup.
Shipper will contact you when arriving.

Food Delivery Platform`, orderID)

	return subject, body
}

// createRestaurantPaymentStatusChangeMessage creates email subject and body for restaurant payment status notifications
func (s *OrderNotificationService) createRestaurantPaymentStatusChangeMessage(orderID, paymentStatus string) (string, string) {
	subject := fmt.Sprintf("Payment %s - Order %s", strings.ToUpper(paymentStatus), orderID)

	var body string
	switch paymentStatus {
	case PaymentStatusPaid:
		body = fmt.Sprintf(`Order %s: Payment confirmed ✓

Start preparation now.

Food Delivery Platform`, orderID)

	default:
		body = fmt.Sprintf(`Order %s: Payment %s

Wait for confirmation before preparing.

Food Delivery Platform`, orderID, strings.ToUpper(paymentStatus))
	}

	return subject, body
}

// createRestaurantOrderCreatedMessage creates email subject and body for restaurant order creation notifications
func (s *OrderNotificationService) createRestaurantOrderCreatedMessage(orderID, _ /* userID */ string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) (string, string) {
	subject := fmt.Sprintf("New Order Received - Order %s", orderID)

	// Build order items list
	var itemsList strings.Builder
	for i, detail := range orderDetails {
		itemsList.WriteString(fmt.Sprintf("%d. Item (Qty: %d) - $%.2f\n", i+1, detail.Quantity, detail.Price))
	}

	body := fmt.Sprintf(`New Order #%s

Total: $%.2f
Payment: %s
Prep Time: %d min

Items:
%s

Start preparation now.

Food Delivery Platform`, orderID, order.TotalPrice, tracking.PaymentStatus, tracking.EstimatedTime, itemsList.String())

	return subject, body
}
