package categorymodel

import "errors"

var (
	ErrNameRequired          = errors.New("name is required")
	ErrCategoryStatusInvalid = errors.New("status must be in (ACTIVE, INACTIVE, DELETED)")
	ErrCategoryIsDeleted     = errors.New("category is deleted")
)
