package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IDeleteCommentRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*foodmodel.FoodRatings, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleteCommentCommandHandler struct {
	repo IDeleteCommentRepo
}

func NewDeleteCommentCommandHandler(repo IDeleteCommentRepo) *DeleteCommentCommandHandler {
	return &DeleteCommentCommandHandler{repo: repo}
}

func (hdl *DeleteCommentCommandHandler) Execute(ctx context.Context, req foodmodel.FoodDeleteCommentReq) error {
	existFoodRating, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodRatingNotFound) {
			return datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodRatingNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if existFoodRating.Status == sharedModel.StatusDelete {
		return datatype.ErrNotFound.WithError(foodmodel.ErrFoodRatingIsDeleted.Error())
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
