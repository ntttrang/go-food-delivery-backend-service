package cartmodel

import "errors"

var (
	ErrUserIdRequired    = errors.New("user id is required")
	ErrFoodIdRequired    = errors.New("food id is required")
	ErrQuantityInvalid   = errors.New("quantity must be greater than 0")
	ErrCartIdRequired    = errors.New("cart id is required")
	ErrCartNotFound      = errors.New("cart not found")
	ErrCartStatusInvalid = errors.New("status must be in (ACTIVE, UPDATED, DELETED)")
	ErrCartIsProcessed   = errors.New("cart is processed. Can't delete")
	ErrFoodInvalid       = errors.New("food is invalid")
)
