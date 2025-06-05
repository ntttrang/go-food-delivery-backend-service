package service

import (
	"fmt"
	"strings"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// Customer message creation methods (continued)

// createCustomerShipperAssignmentMessage creates email subject and body for customer shipper assignment notifications
func (s *OrderNotificationService) createCustomerShipperAssignmentMessage(orderID, _ /* shipperID */ string) (string, string) {
	subject := fmt.Sprintf("Shipper Assigned - Order %s", orderID)

	body := fmt.Sprintf(`Order %s: Shipper assigned!

Your order is being prepared and will be delivered soon.

Track your order in the app.

Food Delivery Team`, orderID)

	return subject, body
}

// createCustomerPaymentStatusChangeMessage creates email subject and body for customer payment status notifications
func (s *OrderNotificationService) createCustomerPaymentStatusChangeMessage(orderID, paymentStatus string) (string, string) {
	subject := fmt.Sprintf("Payment %s - Order %s", strings.ToUpper(paymentStatus), orderID)

	var body string
	switch paymentStatus {
	case PaymentStatusPaid:
		body = fmt.Sprintf(`Order %s: Payment confirmed ✓

Your order is now being prepared.

Food Delivery Team`, orderID)

	case "pending":
		body = fmt.Sprintf(`Order %s: Payment processing...

Please wait for confirmation.

Food Delivery Team`, orderID)

	case "failed":
		body = fmt.Sprintf(`Order %s: Payment failed ✗

Please retry payment in the app.

Food Delivery Team`, orderID)

	default:
		body = fmt.Sprintf(`Order %s: Payment %s

Check app for details.

Food Delivery Team`, orderID, strings.ToUpper(paymentStatus))
	}

	return subject, body
}

// createCustomerOrderCreatedMessage creates email subject and body for customer order creation notifications
func (s *OrderNotificationService) createCustomerOrderCreatedMessage(orderID, _ /* restaurantID */ string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, _ /* orderDetails */ []ordermodel.OrderDetail) (string, string) {
	subject := fmt.Sprintf("Order Confirmed #%s", orderID)

	body := fmt.Sprintf(`Order #%s confirmed ✓

Total: $%.2f
Delivery: %d min

Track in app.

Food Delivery Team`, orderID, order.TotalPrice, tracking.EstimatedTime)

	return subject, body
}

// createCustomerOrderCancelledMessage creates email subject and body for customer order cancellation notifications
func (s *OrderNotificationService) createCustomerOrderCancelledMessage(orderID, reason string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, _ /* orderDetails */ []ordermodel.OrderDetail) (string, string) {
	subject := fmt.Sprintf("Order Cancelled #%s", orderID)

	var refundInfo string
	if tracking.PaymentStatus == PaymentStatusPaid {
		refundInfo = "Refund will be processed within 3-5 business days."
	} else {
		refundInfo = "No payment was processed."
	}

	body := fmt.Sprintf(`Order #%s cancelled

Reason: %s
Total: $%.2f

%s

You can place a new order anytime.

Food Delivery Team`, orderID, reason, order.TotalPrice, refundInfo)

	return subject, body
}
