package service

import (
	"context"
	"fmt"
	"log"
)

// NotificationMessage represents a notification to be sent
type NotificationMessage struct {
	Type      string                 `json:"type"`
	Recipient string                 `json:"recipient"`
	Subject   string                 `json:"subject"`
	Body      string                 `json:"body"`
	Data      map[string]interface{} `json:"data"`
}

// NotificationChannel represents different notification channels
type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelSMS   NotificationChannel = "sms"
	ChannelPush  NotificationChannel = "push"
)

// Repository interfaces for getting user contact information
type IUserNotificationRepo interface {
	GetUserEmail(ctx context.Context, userID string) (string, error)
	GetUserPhone(ctx context.Context, userID string) (string, error)
	GetShipperEmail(ctx context.Context, shipperID string) (string, error)
	GetShipperPhone(ctx context.Context, shipperID string) (string, error)
	GetRestaurantEmail(ctx context.Context, restaurantID string) (string, error)
}

// External notification services interfaces
type IEmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type ISMSService interface {
	SendSMS(ctx context.Context, to, message string) error
}

type IPushNotificationService interface {
	SendPushNotification(ctx context.Context, userID, title, body string, data map[string]interface{}) error
}

// Service
type OrderNotificationService struct {
	userRepo    IUserNotificationRepo
	emailSvc    IEmailService
	smsSvc      ISMSService
	pushSvc     IPushNotificationService
	enabled     bool
}

func NewOrderNotificationService(
	userRepo IUserNotificationRepo,
	emailSvc IEmailService,
	smsSvc ISMSService,
	pushSvc IPushNotificationService,
) *OrderNotificationService {
	return &OrderNotificationService{
		userRepo: userRepo,
		emailSvc: emailSvc,
		smsSvc:   smsSvc,
		pushSvc:  pushSvc,
		enabled:  true, // Can be configured via environment variables
	}
}

// NotifyOrderStateChange sends notifications when order state changes
func (s *OrderNotificationService) NotifyOrderStateChange(ctx context.Context, orderID string, oldState string, newState string) error {
	if !s.enabled {
		return nil
	}

	// For now, just log the notification
	// In a real implementation, this would send actual notifications
	log.Printf("Order %s state changed from %s to %s", orderID, oldState, newState)

	// TODO: Implement actual notification logic
	// 1. Get order details to find user ID and restaurant ID
	// 2. Get user contact information
	// 3. Send appropriate notifications based on state change
	// 4. Send notifications to restaurant if needed

	return s.logNotification("order_state_change", map[string]interface{}{
		"orderID":  orderID,
		"oldState": oldState,
		"newState": newState,
	})
}

// NotifyShipperAssignment sends notifications when a shipper is assigned
func (s *OrderNotificationService) NotifyShipperAssignment(ctx context.Context, orderID string, shipperID string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Shipper %s assigned to order %s", shipperID, orderID)

	// TODO: Implement actual notification logic
	// 1. Get order details
	// 2. Send notification to shipper
	// 3. Send notification to customer
	// 4. Send notification to restaurant

	return s.logNotification("shipper_assignment", map[string]interface{}{
		"orderID":   orderID,
		"shipperID": shipperID,
	})
}

// NotifyPaymentStatusChange sends notifications when payment status changes
func (s *OrderNotificationService) NotifyPaymentStatusChange(ctx context.Context, orderID string, paymentStatus string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Payment status for order %s changed to %s", orderID, paymentStatus)

	// TODO: Implement actual notification logic
	// 1. Get order details
	// 2. Send notification to customer
	// 3. Send notification to restaurant if needed

	return s.logNotification("payment_status_change", map[string]interface{}{
		"orderID":       orderID,
		"paymentStatus": paymentStatus,
	})
}

// NotifyOrderCreated sends notifications when a new order is created
func (s *OrderNotificationService) NotifyOrderCreated(ctx context.Context, orderID string, userID string, restaurantID string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("New order %s created by user %s for restaurant %s", orderID, userID, restaurantID)

	// TODO: Implement actual notification logic
	// 1. Send confirmation to customer
	// 2. Send new order notification to restaurant
	// 3. Send notification to available shippers in the area

	return s.logNotification("order_created", map[string]interface{}{
		"orderID":      orderID,
		"userID":       userID,
		"restaurantID": restaurantID,
	})
}

// NotifyOrderCancelled sends notifications when an order is cancelled
func (s *OrderNotificationService) NotifyOrderCancelled(ctx context.Context, orderID string, reason string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Order %s cancelled. Reason: %s", orderID, reason)

	// TODO: Implement actual notification logic
	// 1. Send cancellation notification to customer
	// 2. Send cancellation notification to restaurant
	// 3. Send cancellation notification to shipper if assigned

	return s.logNotification("order_cancelled", map[string]interface{}{
		"orderID": orderID,
		"reason":  reason,
	})
}

// Helper methods for future implementation

// sendEmailNotification sends an email notification
func (s *OrderNotificationService) sendEmailNotification(ctx context.Context, to, subject, body string) error {
	if s.emailSvc == nil {
		return fmt.Errorf("email service not configured")
	}
	return s.emailSvc.SendEmail(ctx, to, subject, body)
}

// sendSMSNotification sends an SMS notification
func (s *OrderNotificationService) sendSMSNotification(ctx context.Context, to, message string) error {
	if s.smsSvc == nil {
		return fmt.Errorf("SMS service not configured")
	}
	return s.smsSvc.SendSMS(ctx, to, message)
}

// sendPushNotification sends a push notification
func (s *OrderNotificationService) sendPushNotification(ctx context.Context, userID, title, body string, data map[string]interface{}) error {
	if s.pushSvc == nil {
		return fmt.Errorf("push notification service not configured")
	}
	return s.pushSvc.SendPushNotification(ctx, userID, title, body, data)
}

// logNotification logs notification for debugging/monitoring
func (s *OrderNotificationService) logNotification(notificationType string, data map[string]interface{}) error {
	log.Printf("Notification [%s]: %+v", notificationType, data)
	return nil
}

// Enable/Disable notifications
func (s *OrderNotificationService) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func (s *OrderNotificationService) IsEnabled() bool {
	return s.enabled
}
