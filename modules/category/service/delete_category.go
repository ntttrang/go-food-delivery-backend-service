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
type CategoryDeleteReq struct {
	Id uuid.UUID
}

// Initilize service
type IDeleteByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleteByIdCommandHandler struct {
	repo IDeleteByIdRepo
}

func NewDeleteByIdCommandHandler(repo IDeleteByIdRepo) *DeleteByIdCommandHandler {
	return &DeleteByIdCommandHandler{repo: repo}
}

// Implement
func (hdl *DeleteByIdCommandHandler) Execute(ctx context.Context, req CategoryDeleteReq) error {
	category, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, categorymodel.ErrCategoryNotFound) {
			return datatype.ErrNotFound.WithDebug(categorymodel.ErrCategoryNotFound.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if category.Status == sharedModel.StatusDelete {
		return datatype.ErrDeleted.WithError(categorymodel.ErrCategoryIsDeleted.Error())
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil

}
