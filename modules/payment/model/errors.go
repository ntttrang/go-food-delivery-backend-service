package model

import (
	"errors"
)

var (
	ErrCardNotFound          = errors.New("card not found")
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrInvalidCardNumber     = errors.New("invalid card number")
	ErrInvalidCardCVV        = errors.New("invalid card CVV")
	ErrInvalidExpiryDate     = errors.New("invalid expiry date")
	ErrInvalidCardholderName = errors.New("invalid cardholder name")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrInvalidPaymentAmount  = errors.New("invalid payment amount")
	ErrPaymentFailed         = errors.New("payment processing failed")
	ErrUserIDRequired        = errors.New("user ID is required")
	ErrOrderIDRequired       = errors.New("order ID is required")
	ErrMethodRequired        = errors.New("payment method is required")
	ErrProviderRequired      = errors.New("payment provider is required")
	ErrCardTypeRequired      = errors.New("card type is required")
)
