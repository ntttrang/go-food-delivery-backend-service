package ordermodel

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderIdRequired    = errors.New("order id is required")
	ErrUserIdRequired     = errors.New("user id is required")
	ErrTotalPriceInvalid  = errors.New("total price must be greater than 0")
	ErrOrderStatusInvalid = errors.New("order status is invalid")
	ErrOrderIsProcessed   = errors.New("order is already processed")
	ErrOrderIsDelivered   = errors.New("order is already delivered")
	ErrOrderIsCancelled   = errors.New("order is already cancelled")
	ErrRestaurantRequired = errors.New("restaurant id is required")
)
