package foodmodel

import "errors"

var (
	ErrNameRequired           = errors.New("name is required")
	ErrFoodStatusInvalid      = errors.New("status must be in (ACTIVE, INACTIVE, DELETED)")
	ErrFoodIsDeleted          = errors.New("food is deleted")
	ErrFoodNotFound           = errors.New("food not found")
	ErrFieldRequired          = errors.New("food id or user id is required")
	ErrFoodIdRequired         = errors.New("food id is required")
	ErrPointOrCommentRequired = errors.New("point or comment is required")
	ErrPointInvalid           = errors.New("point is 0 to 5")
	ErrUserIdRequired         = errors.New("userId is required")
	ErrFoodLikeNotFound       = errors.New("food like not found")
	ErrFoodRatingNotFound     = errors.New("food ratings not found")
	ErrFoodRatingIsDeleted    = errors.New("food ratings is deleted")
	ErrRestaurantIdEmpty      = errors.New("restaurant Id is empty")
)
