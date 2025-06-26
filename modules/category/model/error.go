package categorymodel

import "errors"

var (
	ErrNameRequired          = errors.New("name is required")
	ErrCategoryStatusInvalid = errors.New("status must be in (ACTIVE, INACTIVE, DELETED)")
	ErrCategoryIsDeleted     = errors.New("category is deleted")
	ErrCategoryNotFound      = errors.New("restaurant not found")
	ErrIdRequired            = errors.New("id is required")
	ErrRequesterRequired     = errors.New("requester information required")
	ErrPermission            = errors.New("only admin and user roles can manage categories")
)
