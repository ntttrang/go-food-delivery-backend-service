package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
)

// Restaurant notification methods

// notifyRestaurantStateChange sends notification to restaurant about order state change
func (s *OrderNotificationService) notifyRestaurantStateChange(ctx context.Context, restaurantID, orderID, oldState, newState string) error {
	// Get userId from restaurantId
	// TODO:
	restaurantIDUuid, _ := uuid.Parse(restaurantID)
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantIDUuid})
	if err != nil {
		return fmt.Errorf("failed to get restaurant: %w", err)
	}
	userId := restaurantMap[restaurantIDUuid].OwnerId

	// Get restaurant email
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userId})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	email := userMap[userId].Email

	// Create notification message
	subject, body := s.createRestaurantStateChangeMessage(orderID, oldState, newState)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to restaurant: %w", err)
	}

	return nil
}

// notifyRestaurantShipperAssignment sends notification to restaurant about shipper assignment
func (s *OrderNotificationService) notifyRestaurantShipperAssignment(ctx context.Context, restaurantID, orderID, shipperID string) error {
	// Get restaurant owner email
	restaurantIDUuid, _ := uuid.Parse(restaurantID)
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantIDUuid})
	if err != nil {
		return fmt.Errorf("failed to get restaurant: %w", err)
	}
	userId := restaurantMap[restaurantIDUuid].OwnerId

	// Get restaurant owner email
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userId})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	email := userMap[userId].Email

	// Create notification message
	subject, body := s.createRestaurantShipperAssignmentMessage(orderID, shipperID)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to restaurant: %w", err)
	}

	return nil
}

// notifyRestaurantPaymentStatusChange sends notification to restaurant about payment status change
func (s *OrderNotificationService) notifyRestaurantPaymentStatusChange(ctx context.Context, restaurantID, orderID, paymentStatus string) error {
	// Get restaurant owner email
	restaurantIDUuid, _ := uuid.Parse(restaurantID)
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantIDUuid})
	if err != nil {
		return fmt.Errorf("failed to get restaurant: %w", err)
	}
	userId := restaurantMap[restaurantIDUuid].OwnerId

	// Get restaurant owner email
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{userId})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	email := userMap[userId].Email

	// Create notification message
	subject, body := s.createRestaurantPaymentStatusChangeMessage(orderID, paymentStatus)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to restaurant: %w", err)
	}

	return nil
}

// notifyRestaurantOrderCreated sends notification to restaurant about new order
func (s *OrderNotificationService) notifyRestaurantOrderCreated(ctx context.Context, restaurantID, orderID, userID string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error {
	// Get restaurant owner email
	restaurantIDUuid, _ := uuid.Parse(restaurantID)
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantIDUuid})
	if err != nil {
		return fmt.Errorf("failed to get restaurant: %w", err)
	}
	ownerUserId := restaurantMap[restaurantIDUuid].OwnerId

	// Get restaurant owner email
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{ownerUserId})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	email := userMap[ownerUserId].Email

	// Create notification message
	subject, body := s.createRestaurantOrderCreatedMessage(orderID, userID, order, tracking, orderDetails)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to restaurant: %w", err)
	}

	return nil
}

// notifyRestaurantOrderCancelled sends cancellation notification to restaurant
func (s *OrderNotificationService) notifyRestaurantOrderCancelled(ctx context.Context, restaurantID, orderID, reason string, order *ordermodel.Order, tracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error {
	// Get restaurant owner email
	restaurantIDUuid, _ := uuid.Parse(restaurantID)
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantIDUuid})
	if err != nil {
		return fmt.Errorf("failed to get restaurant: %w", err)
	}
	ownerUserId := restaurantMap[restaurantIDUuid].OwnerId

	// Get restaurant owner email
	userMap, err := s.userRepo.FindByIds(ctx, []uuid.UUID{ownerUserId})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	email := userMap[ownerUserId].Email

	// Create notification message
	subject, body := s.createRestaurantOrderCancelledMessage(orderID, reason, order, tracking, orderDetails)

	// Send email notification
	if err := s.sendEmailNotification(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to send email to restaurant: %w", err)
	}

	return nil
}
