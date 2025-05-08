package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type CategoryDetailReq struct {
	Id uuid.UUID
}

type CategoryDetailRes struct {
	categorymodel.Category
}

// Initilize service
type IGetDetailRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailRepo
}

func NewGetDetailQueryHandler(repo IGetDetailRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

// Implement
func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req CategoryDetailReq) (*CategoryDetailRes, error) {
	category, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, categorymodel.ErrCategoryNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(categorymodel.ErrCategoryNotFound.Error())
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if category.Status == sharedModel.StatusDelete {
		return nil, datatype.ErrDeleted.WithError(categorymodel.ErrCategoryIsDeleted.Error())
	}

	return &CategoryDetailRes{Category: category}, nil
}
