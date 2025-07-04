package ordermodel

import "errors"

var (
	ErrOrderNotFound             = errors.New("order not found")
	ErrOrderIdRequired           = errors.New("order id is required")
	ErrUserIdRequired            = errors.New("user id is required")
	ErrTotalPriceInvalid         = errors.New("total price must be greater than 0")
	ErrOrderStatusInvalid        = errors.New("order status is invalid")
	ErrOrderIsProcessed          = errors.New("order is already processed")
	ErrOrderIsDelivered          = errors.New("order is already delivered")
	ErrOrderIsCancelled          = errors.New("order is already cancelled")
	ErrRestaurantRequired        = errors.New("restaurant id is required")
	ErrPaymentMethodRequired     = errors.New("payment method is required")
	ErrCardIdRequired            = errors.New("card id is required for card payments")
	ErrDeliveryAddressRequired   = errors.New("delivery address is required")
	ErrCartIdRequired            = errors.New("cart id is required")
	ErrCartNotFound              = errors.New("cart not found")
	ErrCartEmpty                 = errors.New("cart is empty")
	ErrCartAlreadyProcessed      = errors.New("cart has already been processed")
	ErrInvalidPaymentMethod      = errors.New("invalid payment method")
	ErrPaymentFailed             = errors.New("payment processing failed")
	ErrInventoryInsufficient     = errors.New("insufficient inventory")
	ErrRestaurantNotAvailable    = errors.New("restaurant is not available")
	ErrFoodNotAvailable          = errors.New("food item is not available")
	ErrInvalidOrderState         = errors.New("invalid order state transition")
	ErrShipperRequired           = errors.New("shipper id is required")
	ErrMixedRestaurantItems      = errors.New("all cart items must be from the same restaurant")
	ErrInvalidRestaurantIdFormat = errors.New("invalid restaurant ID format")
	ErrInvalidFoodIdFormat       = errors.New("invalid food ID format")
)
