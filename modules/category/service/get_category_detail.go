package service

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IGetDetailRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailRepo
}

func NewGetDetailQueryHandler(repo IGetDetailRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req categorymodel.CategoryDetailReq) (categorymodel.CategoryDetailRes, error) {
	category, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		return categorymodel.CategoryDetailRes{}, err
	}

	if category.Status == sharedModel.StatusDelete {
		return categorymodel.CategoryDetailRes{}, sharedModel.ErrRecordNotFound
	}

	return categorymodel.CategoryDetailRes{Category: category}, nil
}
