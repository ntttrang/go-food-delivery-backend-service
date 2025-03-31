package restaurantmodel

import "errors"

var (
	ErrNameRequired             = errors.New("name is required")
	ErrRestaurantIsDeleted      = errors.New("restaurant is deleted")
	ErrRestaurantNotFound       = errors.New("restaurant not found")
	ErrRestaurantIdRequired     = errors.New("restaurantId is required")
	ErrUserIdRequired           = errors.New("userId is required")
	ErrRestaurantLikeIsDeleted  = errors.New("restaurant like is deleted")
	ErrPointOrCommentRequired   = errors.New("point or comment is required")
	ErrPointInvalid             = errors.New("point is 0 to 5")
	ErrRestaurantRatingNotFound = errors.New("restaurant rating not found")
)
