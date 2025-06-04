package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/rpc-client"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
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

// Repository interfaces for getting order details
type IOrderNotificationRepo interface {
	FindById(ctx context.Context, id string) (*ordermodel.Order, *ordermodel.OrderTracking, []ordermodel.OrderDetail, error)
}

// Repository interfaces for getting user contact information
type IUserNotificationRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]ordermodel.User, error)
}

type IRestaurantNotificationRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcclient.RPCGetByIdsResponseDTO, error)
}

// External notification services interfaces
type IEmailService interface {
	SendEmail(message sharedModel.EmailMessage) error
}

type ISMSService interface {
	SendSMS(ctx context.Context, to, message string) error
}

type IPushNotificationService interface {
	SendPushNotification(ctx context.Context, userID, title, body string, data map[string]interface{}) error
}

// Service
type OrderNotificationService struct {
	orderRepo      IOrderNotificationRepo
	userRepo       IUserNotificationRepo
	restaurantRepo IRestaurantNotificationRepo
	emailSvc       IEmailService
	smsSvc         ISMSService
	pushSvc        IPushNotificationService
	enabled        bool
}

// TODO: System can notify via email, sms (TBD) and push notification(TBD)
func NewOrderNotificationService(
	orderRepo IOrderNotificationRepo,
	userRepo IUserNotificationRepo,
	restaurantRepo IRestaurantNotificationRepo,
	emailSvc IEmailService,
	smsSvc ISMSService,
	pushSvc IPushNotificationService,
) *OrderNotificationService {
	return &OrderNotificationService{
		orderRepo:      orderRepo,
		userRepo:       userRepo,
		restaurantRepo: restaurantRepo,
		emailSvc:       emailSvc,
		smsSvc:         smsSvc,
		pushSvc:        pushSvc,
		enabled:        true, // Can be configured via environment variables
	}
}

// NotifyOrderStateChange sends notifications when order state changes
func (s *OrderNotificationService) NotifyOrderStateChange(ctx context.Context, orderID string, oldState string, newState string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Order %s state changed from %s to %s", orderID, oldState, newState)

	// Get order details to find user ID and restaurant ID
	order, tracking, _, err := s.orderRepo.FindById(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order details for notification: %v", err)
		return err
	}

	// Send notification to customer
	if err := s.notifyCustomerStateChange(ctx, order.UserID, orderID, oldState, newState); err != nil {
		log.Printf("Failed to notify customer about order state change: %v", err)
		// Don't return error, continue with other notifications
	}

	// Send notification to restaurant
	if err := s.notifyRestaurantStateChange(ctx, tracking.RestaurantID, orderID, oldState, newState); err != nil {
		log.Printf("Failed to notify restaurant about order state change: %v", err)
		// Don't return error, continue with other notifications
	}

	// Send notification to shipper if assigned and relevant state
	if order.ShipperID != nil && s.shouldNotifyShipper(newState) {
		if err := s.notifyShipperStateChange(ctx, *order.ShipperID, orderID, oldState, newState); err != nil {
			log.Printf("Failed to notify shipper about order state change: %v", err)
			// Don't return error
		}
	}

	return s.logNotification("order_state_change", map[string]interface{}{
		"orderID":      orderID,
		"oldState":     oldState,
		"newState":     newState,
		"userID":       order.UserID,
		"restaurantID": tracking.RestaurantID,
		"shipperID":    order.ShipperID,
	})
}

// NotifyShipperAssignment sends notifications when a shipper is assigned
func (s *OrderNotificationService) NotifyShipperAssignment(ctx context.Context, orderID string, shipperID string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Shipper %s assigned to order %s", shipperID, orderID)

	// Get order details to find user ID and restaurant ID
	order, tracking, _, err := s.orderRepo.FindById(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order details for shipper assignment notification: %v", err)
		return err
	}

	// Send notification to shipper about the assignment
	if err := s.notifyShipperAssignment(ctx, shipperID, orderID, tracking.RestaurantID); err != nil {
		log.Printf("Failed to notify shipper about assignment: %v", err)
		// Don't return error, continue with other notifications
	}

	// Send notification to customer about shipper assignment
	if err := s.notifyCustomerShipperAssignment(ctx, order.UserID, orderID, shipperID); err != nil {
		log.Printf("Failed to notify customer about shipper assignment: %v", err)
		// Don't return error, continue with other notifications
	}

	// Send notification to restaurant about shipper assignment
	if err := s.notifyRestaurantShipperAssignment(ctx, tracking.RestaurantID, orderID, shipperID); err != nil {
		log.Printf("Failed to notify restaurant about shipper assignment: %v", err)
		// Don't return error
	}

	return s.logNotification("shipper_assignment", map[string]interface{}{
		"orderID":      orderID,
		"shipperID":    shipperID,
		"userID":       order.UserID,
		"restaurantID": tracking.RestaurantID,
	})
}

// NotifyPaymentStatusChange sends notifications when payment status changes
func (s *OrderNotificationService) NotifyPaymentStatusChange(ctx context.Context, orderID string, paymentStatus string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Payment status for order %s changed to %s", orderID, paymentStatus)

	// Get order details to find user ID and restaurant ID
	order, tracking, _, err := s.orderRepo.FindById(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order details for payment status notification: %v", err)
		return err
	}

	// Send notification to customer about payment status change
	if err := s.notifyCustomerPaymentStatusChange(ctx, order.UserID, orderID, paymentStatus); err != nil {
		log.Printf("Failed to notify customer about payment status change: %v", err)
		// Don't return error, continue with other notifications
	}

	// Send notification to restaurant if payment is successful (they can start preparing)
	if paymentStatus == "paid" {
		if err := s.notifyRestaurantPaymentStatusChange(ctx, tracking.RestaurantID, orderID, paymentStatus); err != nil {
			log.Printf("Failed to notify restaurant about payment status change: %v", err)
			// Don't return error
		}
	}

	// Send notification to shipper if assigned and payment is successful
	if order.ShipperID != nil && paymentStatus == "paid" {
		if err := s.notifyShipperPaymentStatusChange(ctx, *order.ShipperID, orderID, paymentStatus); err != nil {
			log.Printf("Failed to notify shipper about payment status change: %v", err)
			// Don't return error
		}
	}

	return s.logNotification("payment_status_change", map[string]interface{}{
		"orderID":       orderID,
		"paymentStatus": paymentStatus,
		"userID":        order.UserID,
		"restaurantID":  tracking.RestaurantID,
		"shipperID":     order.ShipperID,
	})
}

// NotifyOrderCreated sends notifications when a new order is created
func (s *OrderNotificationService) NotifyOrderCreated(ctx context.Context, orderID string, userID string, restaurantID string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("New order %s created by user %s for restaurant %s", orderID, userID, restaurantID)

	// Get order details for comprehensive notification
	order, tracking, orderDetails, err := s.orderRepo.FindById(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order details for order creation notification: %v", err)
		return err
	}

	// Send confirmation to customer
	if err := s.notifyCustomerOrderCreated(ctx, userID, orderID, restaurantID, order, tracking, orderDetails); err != nil {
		log.Printf("Failed to notify customer about order creation: %v", err)
	}

	// Send new order notification to restaurant
	if err := s.notifyRestaurantOrderCreated(ctx, restaurantID, orderID, userID, order, tracking, orderDetails); err != nil {
		log.Printf("Failed to notify restaurant about new order: %v", err)
	}

	// Send notification to available shippers in the area (if shipper is already assigned)
	if order.ShipperID != nil {
		if err := s.notifyShipperOrderCreated(ctx, *order.ShipperID, orderID, restaurantID, order, tracking); err != nil {
			log.Printf("Failed to notify shipper about new order: %v", err)
		}
	}

	return s.logNotification("order_created", map[string]interface{}{
		"orderID":      orderID,
		"userID":       userID,
		"restaurantID": restaurantID,
		"totalPrice":   order.TotalPrice,
		"shipperID":    order.ShipperID,
		"itemCount":    len(orderDetails),
	})
}

// NotifyOrderCancelled sends notifications when an order is cancelled
func (s *OrderNotificationService) NotifyOrderCancelled(ctx context.Context, orderID string, reason string) error {
	if !s.enabled {
		return nil
	}

	log.Printf("Order %s cancelled. Reason: %s", orderID, reason)

	// Get order details to find user ID, restaurant ID, and shipper ID
	order, tracking, orderDetails, err := s.orderRepo.FindById(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order details for cancellation notification: %v", err)
		return err
	}

	// Send cancellation notification to customer
	if err := s.notifyCustomerOrderCancelled(ctx, order.UserID, orderID, reason, order, tracking, orderDetails); err != nil {
		log.Printf("Failed to notify customer about order cancellation: %v", err)
	}

	// Send cancellation notification to restaurant
	if err := s.notifyRestaurantOrderCancelled(ctx, tracking.RestaurantID, orderID, reason, order, tracking, orderDetails); err != nil {
		log.Printf("Failed to notify restaurant about order cancellation: %v", err)
	}

	// Send cancellation notification to shipper if assigned
	if order.ShipperID != nil {
		if err := s.notifyShipperOrderCancelled(ctx, *order.ShipperID, orderID, reason, order, tracking); err != nil {
			log.Printf("Failed to notify shipper about order cancellation: %v", err)
		}
	}

	return s.logNotification("order_cancelled", map[string]interface{}{
		"orderID":      orderID,
		"reason":       reason,
		"userID":       order.UserID,
		"restaurantID": tracking.RestaurantID,
		"shipperID":    order.ShipperID,
		"totalPrice":   order.TotalPrice,
		"itemCount":    len(orderDetails),
	})
}

// Helper methods for state change notifications
// shouldNotifyShipper determines if shipper should be notified for this state
func (s *OrderNotificationService) shouldNotifyShipper(newState string) bool {
	// Notify shipper for states where they are involved
	switch newState {
	case "preparing", "on_the_way", "delivered", "cancel":
		return true
	default:
		return false
	}
}

// notifyRestaurantStateChange sends notification to restaurant about order state change
// sendEmailNotification sends an email notification
func (s *OrderNotificationService) sendEmailNotification(ctx context.Context, to, subject, body string) error {
	if s.emailSvc == nil {
		return fmt.Errorf("email service not configured")
	}
	var msg sharedModel.EmailMessage
	msg.From = "sender@gmail.com"
	msg.To = []string{to}
	msg.Subject = subject
	msg.Body = body
	return s.emailSvc.SendEmail(msg)
}

// sendSMSNotification sends an SMS notification
// TODO: [TBD]
func (s *OrderNotificationService) sendSMSNotification(ctx context.Context, to, message string) error {
	if s.smsSvc == nil {
		return fmt.Errorf("SMS service not configured")
	}
	return s.smsSvc.SendSMS(ctx, to, message)
}

// sendPushNotification sends a push notification
// TODO: [TBD]
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

// Message creation helper methods

// getStateDisplayName returns a user-friendly display name for order states
func (s *OrderNotificationService) getStateDisplayName(state string) string {
	switch state {
	case "waiting_for_shipper":
		return "waiting for shipper"
	case "preparing":
		return "being prepared"
	case "on_the_way":
		return "on the way"
	case "delivered":
		return "delivered"
	case "cancel":
		return "cancelled"
	default:
		return state
	}
}
