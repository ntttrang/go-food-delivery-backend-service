package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type FoodDetailReq struct {
	Id uuid.UUID
}

type FoodDetailRes struct {
	foodmodel.Food
}

// Initilize service
type IGetDetailRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailRepo
}

func NewGetDetailQueryHandler(repo IGetDetailRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

// Implement
func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req FoodDetailReq) (*FoodDetailRes, error) {
	food, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodNotFound.Error())
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if food.Status == string(datatype.StatusDeleted) {
		return nil, datatype.ErrDeleted.WithError(foodmodel.ErrFoodIsDeleted.Error())
	}

	return &FoodDetailRes{Food: food}, nil
}
