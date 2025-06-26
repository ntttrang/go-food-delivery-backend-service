package usermodel

import "errors"

var (
	ErrEmailRequired       = errors.New("email is required")
	ErrEmailInvalid        = errors.New("email is invalid")
	ErrPasswordInvalid     = errors.New("password must be greater than 6 characters")
	ErrPasswordRequired    = errors.New("password is required")
	ErrFirstNameRequired   = errors.New("first name is required")
	ErrLastNameRequired    = errors.New("last name is required")
	ErrUserDeletedOrBanned = errors.New("user is deleted or banned")
	ErrUserNotFound        = errors.New("user not found")
	ErrIdRequired          = errors.New("id is required")
	ErrInvalidPhoneNumber  = errors.New("phone number invalid")
	ErrAddrRequired        = errors.New("addr is required")
	ErrDuplicated          = errors.New("addr is duplicated")
	ErrUserAddrNotFound    = errors.New("user address not found")
	ErrPermission          = errors.New("you can only update your own profile unless you're an admin")
)
