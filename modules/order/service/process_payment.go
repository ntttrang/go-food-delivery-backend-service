package service

import (
	"context"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// PaymentCard represents a payment card
type PaymentCard struct {
	ID             uuid.UUID `json:"id"`
	Method         string    `json:"method"`
	Provider       string    `json:"provider"`
	CardholderName string    `json:"cardholderName"`
	CardNumber     string    `json:"cardNumber"`
	CardType       string    `json:"cardType"`
	ExpiryMonth    string    `json:"expiryMonth"`
	ExpiryYear     string    `json:"expiryYear"`
	UserID         uuid.UUID `json:"userId"`
	Status         string    `json:"status"`
}

// PaymentRequest represents a payment processing request
type PaymentRequest struct {
	OrderID       string  `json:"orderId"`
	UserID        string  `json:"userId"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	CardID        *string `json:"cardId,omitempty"`
}

// PaymentResult represents the result of payment processing
type PaymentResult struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transactionId,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
}

// Repository interfaces
type IPaymentCardRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*ordermodel.Card, error)
}

// Service
type PaymentProcessingService struct {
	cardRepo IPaymentCardRepo
}

func NewPaymentProcessingService(
	cardRepo IPaymentCardRepo,
) *PaymentProcessingService {
	return &PaymentProcessingService{
		cardRepo: cardRepo,
	}
}

// ValidatePaymentMethod validates the payment method and associated data
func (s *PaymentProcessingService) ValidatePaymentMethod(ctx context.Context, req *PaymentRequest) error {
	// Validate payment method
	if req.PaymentMethod != MethodCash && req.PaymentMethod != MethodCreditCard && req.PaymentMethod != MethodDebitCard {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidPaymentMethod).WithDebug("payment method must be 'CASH' or 'CREDIT_CARD' or 'DEBIT_CARD' ")
	}

	// For card payments, validate card exists and belongs to user
	if req.PaymentMethod == MethodCreditCard || req.PaymentMethod == MethodDebitCard {
		if req.CardID == nil {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrCardIdRequired).WithDebug(ordermodel.ErrCardIdRequired.Error())
		}

		cardID, err := uuid.Parse(*req.CardID)
		if err != nil {
			return datatype.ErrBadRequest.WithError("invalid card ID format")
		}

		card, err := s.cardRepo.FindById(ctx, cardID)
		if err != nil {
			return datatype.ErrNotFound.WithWrap(err).WithDebug("card not found")
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return datatype.ErrBadRequest.WithError("invalid user ID format")
		}

		if card.UserID != userID {
			return datatype.ErrForbidden.WithError("card does not belong to user")
		}

		if card.Status != "ACTIVE" {
			return datatype.ErrBadRequest.WithError("card is not active")
		}
	}

	return nil
}

// ProcessPayment processes the payment based on the method
func (s *PaymentProcessingService) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	// Validate payment method first
	if err := s.ValidatePaymentMethod(ctx, req); err != nil {
		return nil, err
	}

	switch req.PaymentMethod {
	case MethodCash:
		return s.processCashPayment(ctx, req)
	case MethodDebitCard, MethodCreditCard:
		return s.processCardPayment(ctx, req)
	default:
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidPaymentMethod).WithDebug("unsupported payment method")
	}
}

// processCashPayment handles cash payment
func (s *PaymentProcessingService) processCashPayment(_ context.Context, req *PaymentRequest) (*PaymentResult, error) {
	// For cash payments, we just mark as pending
	// The payment will be completed when the order is delivered

	// Default implementation - simulate successful payment
	return &PaymentResult{
		Success:       true,
		TransactionID: req.PaymentMethod + "_" + req.OrderID,
	}, nil
}

// processCardPayment handles card payment processing
func (s *PaymentProcessingService) processCardPayment(_ context.Context, req *PaymentRequest) (*PaymentResult, error) {
	if req.CardID == nil {
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrCardIdRequired).WithDebug(ordermodel.ErrCardIdRequired.Error())
	}

	cardID, err := uuid.Parse(*req.CardID)
	if err != nil {
		return nil, datatype.ErrBadRequest.WithError("invalid card ID format")
	}

	// TODO: (TBD)
	// In a real implementation, this would integrate with Stripe, PayPal, etc.
	// Step 1: Get card info
	// Step 2: Process payment through gateway

	// Default implementation - simulate successful payment
	return &PaymentResult{
		Success:       true,
		TransactionID: req.PaymentMethod + "_" + req.OrderID + "_" + cardID.String(),
	}, nil
}

// GetPaymentStatus determines the payment status based on payment method and result
// func (s *PaymentProcessingService) GetPaymentStatus(paymentMethod string, result *PaymentResult) string {
// 	if !result.Success {
// 		return "failed"
// 	}

// 	switch paymentMethod {
// 	case MethodCash:
// 		return "pending" // Cash payments are pending until delivery
// 	case MethodDebitCard, MethodCreditCard:
// 		return PaymentStatusPaid // Card payments are immediately paid if successful
// 	default:
// 		return "pending"
// 	}
// }
