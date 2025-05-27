package service

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type FoodCommentCreateReq struct {
	UserId  uuid.UUID `json:"-"` // Get from Token
	FoodId  uuid.UUID `json:"foodId"`
	Point   *float64  `json:"point"`
	Comment *string   `json:"comment"`

	Id uuid.UUID `json:"-"` // Internal BE
}

func (r FoodCommentCreateReq) Validate() error {
	if r.UserId.String() == "" {
		return foodmodel.ErrUserIdRequired
	}
	if r.FoodId.String() == "" {
		return foodmodel.ErrFoodIdRequired
	}
	if r.Point == nil && (r.Comment == nil || *r.Comment == "") {
		return foodmodel.ErrPointOrCommentRequired
	}
	if r.Point != nil && *r.Point > 5.0 {
		return foodmodel.ErrPointInvalid
	}
	return nil
}

func (FoodCommentCreateReq) TableName() string {
	return foodmodel.FoodRatings{}.TableName()
}

// Initilize service
type IInsertCommentFoodRepo interface {
	Insert(ctx context.Context, req *FoodCommentCreateReq) error
}

type CreateFoodCommentCommandHandler struct {
	repo IInsertCommentFoodRepo
}

func NewCommentFoodCommandHandler(repo IInsertCommentFoodRepo) *CreateFoodCommentCommandHandler {
	return &CreateFoodCommentCommandHandler{repo: repo}
}

// Implement
func (hdl *CreateFoodCommentCommandHandler) Execute(ctx context.Context, req *FoodCommentCreateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	req.Id, _ = uuid.NewV7()
	if err := hdl.repo.Insert(ctx, req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}
