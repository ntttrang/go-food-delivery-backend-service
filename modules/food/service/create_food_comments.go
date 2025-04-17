package service

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IInsertCommentFoodRepo interface {
	Insert(ctx context.Context, req *foodmodel.FoodCommentCreateReq) error
}

type CreateFoodCommentCommandHandler struct {
	repo IInsertCommentFoodRepo
}

func NewCommentFoodCommandHandler(repo IInsertCommentFoodRepo) *CreateFoodCommentCommandHandler {
	return &CreateFoodCommentCommandHandler{repo: repo}
}

func (hdl *CreateFoodCommentCommandHandler) Execute(ctx context.Context, req *foodmodel.FoodCommentCreateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	req.Id, _ = uuid.NewV7()
	if err := hdl.repo.Insert(ctx, req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}
