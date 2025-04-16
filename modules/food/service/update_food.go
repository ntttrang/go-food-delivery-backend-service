package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IUpdateByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
	Update(ctx context.Context, id uuid.UUID, dto foodmodel.FoodUpdateReq) error
}

type UpdateCommandHandler struct {
	repo IUpdateByIdRepo
}

func NewUpdateCommandHandler(repo IUpdateByIdRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: repo}
}

func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req foodmodel.FoodUpdateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	category, err := hdl.repo.FindById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodNotFound) {
			return datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodNotFound.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if category.Status == sharedModel.StatusDelete {
		return datatype.ErrDeleted.WithError(foodmodel.ErrFoodIsDeleted.Error())
	}

	if err := hdl.repo.Update(ctx, req.Id, req); err != nil {
		return err
	}

	return nil
}
