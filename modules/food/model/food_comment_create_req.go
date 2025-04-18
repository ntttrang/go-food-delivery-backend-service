package foodmodel

import "github.com/google/uuid"

type FoodCommentCreateReq struct {
	UserId  uuid.UUID `json:"-"` // Get from Token
	FoodId  uuid.UUID `json:"restaurantId"`
	Point   *float64  `json:"point"`
	Comment *string   `json:"comment"`

	Id uuid.UUID `json:"-"` // Internal BE
}

func (r FoodCommentCreateReq) Validate() error {
	if r.UserId.String() == "" {
		return ErrUserIdRequired
	}
	if r.FoodId.String() == "" {
		return ErrFoodIdRequired
	}
	if r.Point == nil && (r.Comment == nil || *r.Comment == "") {
		return ErrPointOrCommentRequired
	}
	if r.Point != nil && *r.Point > 5.0 {
		return ErrPointInvalid
	}
	return nil
}

func (FoodCommentCreateReq) TableName() string {
	return FoodRatings{}.TableName()
}
