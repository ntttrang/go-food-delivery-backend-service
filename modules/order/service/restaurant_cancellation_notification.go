package service

import (
	"fmt"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// createRestaurantOrderCancelledMessage creates email subject and body for restaurant order cancellation notifications
func (s *OrderNotificationService) createRestaurantOrderCancelledMessage(orderID, reason string, order *ordermodel.Order, _ /* tracking */ *ordermodel.OrderTracking, _ /* orderDetails */ []ordermodel.OrderDetail) (string, string) {
	subject := fmt.Sprintf("Order Cancelled #%s", orderID)

	body := fmt.Sprintf(`Order #%s cancelled

Reason: %s
Total: $%.2f

Stop preparation if started.
Continue accepting new orders.

Food Delivery Platform`, orderID, reason, order.TotalPrice)

	return subject, body
}
