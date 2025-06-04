package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OrderEvent represents an order-related event
type OrderEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	OrderID     string                 `json:"orderId"`
	UserID      string                 `json:"userId,omitempty"`
	Data        map[string]interface{} `json:"data"`
	OccurredAt  time.Time              `json:"occurredAt"`
	Version     string                 `json:"version"`
}

// Event types
const (
	OrderCreatedType           = "order.created"
	OrderStateChangedType      = "order.state_changed"
	OrderPaymentProcessedType  = "order.payment_processed"
	OrderShipperAssignedType   = "order.shipper_assigned"
	OrderCancelledType         = "order.cancelled"
	OrderDeliveredType         = "order.delivered"
)

// NewOrderEvent creates a new order event
func NewOrderEvent(eventType, orderID, userID string, data map[string]interface{}) *OrderEvent {
	return &OrderEvent{
		ID:         uuid.New().String(),
		Type:       eventType,
		OrderID:    orderID,
		UserID:     userID,
		Data:       data,
		OccurredAt: time.Now(),
		Version:    "1.0",
	}
}

// ToJSON converts the event to JSON bytes
func (e *OrderEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON creates an OrderEvent from JSON bytes
func FromJSON(data []byte) (*OrderEvent, error) {
	var event OrderEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

// Helper functions to create specific events

// NewOrderCreatedEvent creates an order created event
func NewOrderCreatedEvent(orderID, userID, restaurantID string, totalPrice float64, orderDetails []map[string]interface{}) *OrderEvent {
	data := map[string]interface{}{
		"restaurantId":  restaurantID,
		"totalPrice":    totalPrice,
		"orderDetails":  orderDetails,
		"createdAt":     time.Now(),
	}
	return NewOrderEvent(OrderCreatedType, orderID, userID, data)
}

// NewOrderStateChangedEvent creates an order state changed event
func NewOrderStateChangedEvent(orderID, oldState, newState, updatedBy string) *OrderEvent {
	data := map[string]interface{}{
		"oldState":  oldState,
		"newState":  newState,
		"updatedBy": updatedBy,
		"changedAt": time.Now(),
	}
	return NewOrderEvent(OrderStateChangedType, orderID, "", data)
}

// NewOrderPaymentProcessedEvent creates a payment processed event
func NewOrderPaymentProcessedEvent(orderID, paymentStatus, paymentMethod string, amount float64, success bool) *OrderEvent {
	data := map[string]interface{}{
		"paymentStatus": paymentStatus,
		"paymentMethod": paymentMethod,
		"amount":        amount,
		"success":       success,
		"processedAt":   time.Now(),
	}
	return NewOrderEvent(OrderPaymentProcessedType, orderID, "", data)
}

// NewOrderShipperAssignedEvent creates a shipper assigned event
func NewOrderShipperAssignedEvent(orderID, shipperID, assignedBy string) *OrderEvent {
	data := map[string]interface{}{
		"shipperID":  shipperID,
		"assignedBy": assignedBy,
		"assignedAt": time.Now(),
	}
	return NewOrderEvent(OrderShipperAssignedType, orderID, "", data)
}

// NewOrderCancelledEvent creates an order cancelled event
func NewOrderCancelledEvent(orderID, reason, cancelledBy string) *OrderEvent {
	data := map[string]interface{}{
		"reason":      reason,
		"cancelledBy": cancelledBy,
		"cancelledAt": time.Now(),
	}
	return NewOrderEvent(OrderCancelledType, orderID, "", data)
}

// NewOrderDeliveredEvent creates an order delivered event
func NewOrderDeliveredEvent(orderID, shipperID string, deliveryTime int) *OrderEvent {
	data := map[string]interface{}{
		"shipperID":    shipperID,
		"deliveryTime": deliveryTime,
		"deliveredAt":  time.Now(),
	}
	return NewOrderEvent(OrderDeliveredType, orderID, "", data)
}
