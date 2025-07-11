package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type FoodDeleteReq struct {
	Id uuid.UUID
}

// Initilize service
type IDeleteByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleteByIdCommandHandler struct {
	repo IDeleteByIdRepo
}

func NewDeleteByIdCommandHandler(repo IDeleteByIdRepo) *DeleteByIdCommandHandler {
	return &DeleteByIdCommandHandler{repo: repo}
}

// Implement
func (hdl *DeleteByIdCommandHandler) Execute(ctx context.Context, req FoodDeleteReq) error {
	food, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodNotFound) {
			return datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodNotFound.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if food.Status == string(datatype.StatusDeleted) {
		return datatype.ErrDeleted.WithError(foodmodel.ErrFoodIsDeleted.Error())
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
