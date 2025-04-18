package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IFoodLikeRepo interface {
	FindByUserIdAndFoodId(ctx context.Context, userId, restaurantId uuid.UUID) (*foodmodel.FoodLike, error)
	Insert(ctx context.Context, req foodmodel.FoodLike) error
	Delete(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) error
}

type AddFavoritesCommandHandler struct {
	repo IFoodLikeRepo
}

func NewAddFavoritesCommandHandler(repo IFoodLikeRepo) *AddFavoritesCommandHandler {
	return &AddFavoritesCommandHandler{
		repo: repo,
	}
}

func (hdl *AddFavoritesCommandHandler) Execute(ctx context.Context, req foodmodel.FoodLike) (*string, error) {
	if err := req.Validate(); err != nil {
		return nil, datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	restaurantLike, err := hdl.repo.FindByUserIdAndFoodId(ctx, req.UserId, req.FoodId)

	if err != nil {
		if !errors.Is(err, foodmodel.ErrFoodLikeNotFound) {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
	}

	flag := ""
	if restaurantLike == nil {
		if err := hdl.repo.Insert(ctx, req); err != nil {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
		flag = "Add to My favorite food"
	} else {
		if err := hdl.repo.Delete(ctx, req.FoodId, req.UserId); err != nil {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
		flag = "Remove out My favorite food"
	}

	return &flag, nil
}
