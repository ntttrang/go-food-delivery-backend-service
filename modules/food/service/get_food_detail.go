package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IGetDetailRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailRepo
}

func NewGetDetailQueryHandler(repo IGetDetailRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req foodmodel.FoodDetailReq) (foodmodel.FoodDetailRes, error) {
	food, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodNotFound) {
			return foodmodel.FoodDetailRes{}, datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodNotFound.Error())
		}
		return foodmodel.FoodDetailRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if food.Status == sharedModel.StatusDelete {
		return foodmodel.FoodDetailRes{}, datatype.ErrDeleted.WithError(foodmodel.ErrFoodIsDeleted.Error())
	}

	return foodmodel.FoodDetailRes{Food: food}, nil
}
