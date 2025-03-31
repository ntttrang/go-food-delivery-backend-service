package service

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IUpdateByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error)
	Update(ctx context.Context, id uuid.UUID, dto categorymodel.CategoryUpdateReq) error
}

type UpdateCommandHandler struct {
	repo IUpdateByIdRepo
}

func NewUpdateCommandHandler(repo IUpdateByIdRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: repo}
}

func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req categorymodel.CategoryUpdateReq) error {
	if err := req.Validate(); err != nil {
		return err
	}

	category, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		return err
	}

	if category.Status == sharedModel.StatusDelete {
		return categorymodel.ErrCategoryIsDeleted
	}

	if err := hdl.repo.Update(ctx, req.Id, req); err != nil {
		return err
	}

	return nil
}
