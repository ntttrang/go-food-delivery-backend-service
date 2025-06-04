package main

import (
	"log"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/events"
)

func main() {
	// Create Kafka consumer
	brokers := []string{"localhost:9092"}
	topic := "order-events"

	consumer, err := events.NewKafkaConsumer(brokers, topic)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Define event handler
	handler := func(event *events.OrderEvent) error {
		log.Printf("Received event: %s", event.Type)
		log.Printf("Order ID: %s", event.OrderID)
		log.Printf("User ID: %s", event.UserID)
		log.Printf("Data: %+v", event.Data)
		log.Printf("Occurred at: %s", event.OccurredAt)
		
		// Handle different event types
		switch event.Type {
		case events.OrderCreatedType:
			log.Printf("📦 New order created: %s", event.OrderID)
			// TODO: Send notification, update inventory, etc.
			
		case events.OrderStateChangedType:
			oldState := event.Data["oldState"]
			newState := event.Data["newState"]
			log.Printf("📋 Order %s state changed from %s to %s", event.OrderID, oldState, newState)
			// TODO: Update tracking, send notifications, etc.
			
		case events.OrderPaymentProcessedType:
			success := event.Data["success"]
			log.Printf("💳 Payment processed for order %s: success=%v", event.OrderID, success)
			// TODO: Update payment status, send confirmation, etc.
			
		case events.OrderShipperAssignedType:
			shipperID := event.Data["shipperID"]
			log.Printf("🚚 Shipper %s assigned to order %s", shipperID, event.OrderID)
			// TODO: Notify shipper, update tracking, etc.
			
		case events.OrderCancelledType:
			reason := event.Data["reason"]
			log.Printf("❌ Order %s cancelled: %s", event.OrderID, reason)
			// TODO: Refund payment, update inventory, send notification, etc.
			
		case events.OrderDeliveredType:
			log.Printf("✅ Order %s delivered successfully", event.OrderID)
			// TODO: Complete order, send confirmation, update analytics, etc.
		}
		
		return nil
	}

	log.Println("Starting Kafka consumer...")
	log.Println("Listening for order events...")
	
	// Start consuming events (this will block)
	if err := consumer.ConsumeEvents(handler); err != nil {
		log.Fatalf("Failed to consume events: %v", err)
	}
}
