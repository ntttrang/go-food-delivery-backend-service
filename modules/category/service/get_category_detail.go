package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
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
		if errors.Is(err, categorymodel.ErrCategoryNotFound) {
			return categorymodel.CategoryDetailRes{}, datatype.ErrNotFound.WithDebug(categorymodel.ErrCategoryNotFound.Error())
		}
		return categorymodel.CategoryDetailRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if category.Status == sharedModel.StatusDelete {
		return categorymodel.CategoryDetailRes{}, datatype.ErrDeleted.WithError(categorymodel.ErrCategoryIsDeleted.Error())
	}

	return categorymodel.CategoryDetailRes{Category: category}, nil
}
