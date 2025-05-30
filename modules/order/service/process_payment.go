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
	FindByID(ctx context.Context, id uuid.UUID) (*PaymentCard, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]PaymentCard, error)
}

// External payment gateway interface (for future implementation)
type IPaymentGateway interface {
	ProcessCardPayment(ctx context.Context, card *PaymentCard, amount float64, orderID string) (*PaymentResult, error)
	ProcessCashPayment(ctx context.Context, amount float64, orderID string) (*PaymentResult, error)
}

// Service
type PaymentProcessingService struct {
	cardRepo       IPaymentCardRepo
	paymentGateway IPaymentGateway
}

func NewPaymentProcessingService(
	cardRepo IPaymentCardRepo,
	paymentGateway IPaymentGateway,
) *PaymentProcessingService {
	return &PaymentProcessingService{
		cardRepo:       cardRepo,
		paymentGateway: paymentGateway,
	}
}

// ValidatePaymentMethod validates the payment method and associated data
func (s *PaymentProcessingService) ValidatePaymentMethod(ctx context.Context, req *PaymentRequest) error {
	// Validate payment method
	if req.PaymentMethod != "cash" && req.PaymentMethod != "card" {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidPaymentMethod).WithDebug("payment method must be 'cash' or 'card'")
	}

	// For card payments, validate card exists and belongs to user
	if req.PaymentMethod == "card" {
		if req.CardID == nil {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrCardIdRequired).WithDebug(ordermodel.ErrCardIdRequired.Error())
		}

		cardID, err := uuid.Parse(*req.CardID)
		if err != nil {
			return datatype.ErrBadRequest.WithError("invalid card ID format")
		}

		card, err := s.cardRepo.FindByID(ctx, cardID)
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
	case "cash":
		return s.processCashPayment(ctx, req)
	case "card":
		return s.processCardPayment(ctx, req)
	default:
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrInvalidPaymentMethod).WithDebug("unsupported payment method")
	}
}

// processCashPayment handles cash payment (no actual processing needed)
func (s *PaymentProcessingService) processCashPayment(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	// For cash payments, we just mark as pending
	// The payment will be completed when the order is delivered
	if s.paymentGateway != nil {
		return s.paymentGateway.ProcessCashPayment(ctx, req.Amount, req.OrderID)
	}

	// Default implementation - cash payments are always "successful" but pending
	return &PaymentResult{
		Success:       true,
		TransactionID: "cash_" + req.OrderID,
	}, nil
}

// processCardPayment handles card payment processing
func (s *PaymentProcessingService) processCardPayment(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
	if req.CardID == nil {
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrCardIdRequired).WithDebug(ordermodel.ErrCardIdRequired.Error())
	}

	cardID, err := uuid.Parse(*req.CardID)
	if err != nil {
		return nil, datatype.ErrBadRequest.WithError("invalid card ID format")
	}

	card, err := s.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return nil, datatype.ErrNotFound.WithWrap(err).WithDebug("card not found")
	}

	// Process payment through gateway
	if s.paymentGateway != nil {
		result, err := s.paymentGateway.ProcessCardPayment(ctx, card, req.Amount, req.OrderID)
		if err != nil {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug("payment gateway error")
		}
		return result, nil
	}

	// Default implementation - simulate successful payment
	// In a real implementation, this would integrate with Stripe, PayPal, etc.
	return &PaymentResult{
		Success:       true,
		TransactionID: "card_" + req.OrderID + "_" + cardID.String(),
	}, nil
}

// GetPaymentStatus determines the payment status based on payment method and result
func (s *PaymentProcessingService) GetPaymentStatus(paymentMethod string, result *PaymentResult) string {
	if !result.Success {
		return "failed"
	}

	switch paymentMethod {
	case "cash":
		return "pending" // Cash payments are pending until delivery
	case "card":
		return "paid" // Card payments are immediately paid if successful
	default:
		return "pending"
	}
}
