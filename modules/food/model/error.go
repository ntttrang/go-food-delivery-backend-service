package foodmodel

import "errors"

var (
	ErrNameRequired      = errors.New("name is required")
	ErrFoodStatusInvalid = errors.New("status must be in (ACTIVE, INACTIVE, DELETED)")
	ErrFoodIsDeleted     = errors.New("Food is deleted")
	ErrFoodNotFound      = errors.New("food not found")
)
