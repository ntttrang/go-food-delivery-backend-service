package service

import (
	"context"
	"log"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/events"
)

// KafkaEventDrivenOrderService provides event-driven order operations using Kafka
type KafkaEventDrivenOrderService struct {
	createHandler *CreateCommandHandler
	stateHandler  *OrderStateManagementService
	producer      *events.KafkaProducer
}

func NewKafkaEventDrivenOrderService(
	createHandler *CreateCommandHandler,
	stateHandler *OrderStateManagementService,
	producer *events.KafkaProducer,
) *KafkaEventDrivenOrderService {
	return &KafkaEventDrivenOrderService{
		createHandler: createHandler,
		stateHandler:  stateHandler,
		producer:      producer,
	}
}

// CreateOrder creates an order and publishes events to Kafka
func (s *KafkaEventDrivenOrderService) CreateOrder(ctx context.Context, data *OrderCreateDto) (string, error) {
	// Create the order using existing service
	orderID, err := s.createHandler.Execute(ctx, data)
	if err != nil {
		return "", err
	}

	// Create order created event
	orderCreatedEvent := events.NewOrderCreatedEvent(
		orderID,
		data.UserID,
		data.RestaurantID,
		data.TotalPrice,
		s.convertOrderDetailsToEventData(data.OrderDetails),
	)

	// Publish event to Kafka
	if err := s.producer.PublishEvent(orderCreatedEvent); err != nil {
		log.Printf("Failed to publish OrderCreated event to Kafka for order %s: %v", orderID, err)
		// Don't fail the order creation if event publishing fails
	}

	return orderID, nil
}

// Helper method to convert order details to event data
func (s *KafkaEventDrivenOrderService) convertOrderDetailsToEventData(details []OrderDetailCreateDto) []map[string]interface{} {
	var eventDetails []map[string]interface{}
	for _, detail := range details {
		eventDetail := map[string]interface{}{
			"price":    detail.Price,
			"quantity": detail.Quantity,
			"discount": detail.Discount,
		}

		// Add food origin data if available
		if detail.FoodOrigin != nil {
			eventDetail["foodId"] = detail.FoodOrigin.Id
			eventDetail["foodName"] = detail.FoodOrigin.Name
		}

		eventDetails = append(eventDetails, eventDetail)
	}
	return eventDetails
}
