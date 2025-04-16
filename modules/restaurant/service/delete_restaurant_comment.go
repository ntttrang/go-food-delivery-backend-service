package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IDeleteCommentRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.RestaurantRating, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleteCommentCommandHandler struct {
	repo IDeleteCommentRepo
}

func NewDeleteCommentCommandHandler(repo IDeleteCommentRepo) *DeleteCommentCommandHandler {
	return &DeleteCommentCommandHandler{repo: repo}
}

func (hdl *DeleteCommentCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantDeleteCommentReq) error {
	existRestaurantRating, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantRatingNotFound) {
			return datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantRatingNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if existRestaurantRating.Status == sharedModel.StatusDelete {
		return datatype.ErrNotFound.WithError(restaurantmodel.ErrRestaurantRatingIsDeleted.Error())
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
