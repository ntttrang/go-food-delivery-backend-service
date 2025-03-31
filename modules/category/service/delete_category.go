package service

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

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

func (hdl *DeleteByIdCommandHandler) Execute(ctx context.Context, req categorymodel.CategoryDeleteReq) error {
	category, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		return err
	}

	if category.Status == sharedModel.StatusDelete {
		return categorymodel.ErrCategoryIsDeleted
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		return err
	}

	return nil

}
