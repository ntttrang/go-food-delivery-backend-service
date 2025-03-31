package restaurantmodel

import "github.com/google/uuid"

type RestaurantCommentCreateReq struct {
	UserID       string   `json:"userId"`
	RestaurantID string   `json:"restaurantId"`
	Point        *float64 `json:"point"`
	Comment      *string  `json:"comment"`

	Id uuid.UUID `json:"-"` // Internal BE
}

func (r RestaurantCommentCreateReq) Validate() error {
	if r.UserID == "" {
		return ErrUserIdRequired
	}
	if r.RestaurantID == "" {
		return ErrRestaurantIdRequired
	}
	if r.Point == nil && (r.Comment == nil || *r.Comment == "") {
		return ErrPointOrCommentRequired
	}
	if r.Point != nil && *r.Point > 5.0 {
		return ErrPointInvalid
	}
	return nil
}

func (RestaurantCommentCreateReq) TableName() string {
	return RestaurantRating{}.TableName()
}
