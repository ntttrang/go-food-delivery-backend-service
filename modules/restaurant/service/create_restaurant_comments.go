package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type RestaurantCommentCreateReq struct {
	UserID       uuid.UUID `json:"-"` // Get from Token
	RestaurantID uuid.UUID `json:"restaurantId"`
	Point        *float64  `json:"point"`
	Comment      *string   `json:"comment"`

	Id uuid.UUID `json:"-"` // Internal BE
}

func (r RestaurantCommentCreateReq) Validate() error {
	if r.UserID.String() == "" {
		return restaurantmodel.ErrUserIdRequired
	}
	if r.RestaurantID.String() == "" {
		return restaurantmodel.ErrRestaurantIdRequired
	}
	if r.Point == nil && (r.Comment == nil || *r.Comment == "") {
		return restaurantmodel.ErrPointOrCommentRequired
	}
	if r.Point != nil && *r.Point > 5.0 {
		return restaurantmodel.ErrPointInvalid
	}
	return nil
}

func (RestaurantCommentCreateReq) TableName() string {
	return restaurantmodel.RestaurantRating{}.TableName()
}

// Initialize service
type IInsertCommentRestaurantRepo interface {
	Insert(ctx context.Context, req *RestaurantCommentCreateReq) error
}

type CreateRestaurantCommentCommandHandler struct {
	repo IInsertCommentRestaurantRepo
}

func NewCommentRestaurantCommandHandler(repo IInsertCommentRestaurantRepo) *CreateRestaurantCommentCommandHandler {
	return &CreateRestaurantCommentCommandHandler{repo: repo}
}

// Implement
func (hdl *CreateRestaurantCommentCommandHandler) Execute(ctx context.Context, req *RestaurantCommentCreateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	req.Id, _ = uuid.NewV7()
	if err := hdl.repo.Insert(ctx, req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}
