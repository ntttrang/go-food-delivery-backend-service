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
type CategoryUpdateReq struct {
	// Use pointer to accept empty string
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`

	Id uuid.UUID `json:"-"`
}

func (CategoryUpdateReq) TableName() string {
	return categorymodel.Category{}.TableName()
}

func (c CategoryUpdateReq) validate() error {
	if c.Status != nil && *c.Status != sharedModel.StatusActive && *c.Status != sharedModel.StatusDelete && *c.Status != sharedModel.StatusInactive {
		return categorymodel.ErrCategoryStatusInvalid
	}
	return nil
}

// Initilize service
type IUpdateByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (categorymodel.Category, error)
	Update(ctx context.Context, id uuid.UUID, dto CategoryUpdateReq) error
}

type UpdateCommandHandler struct {
	repo IUpdateByIdRepo
}

func NewUpdateCommandHandler(repo IUpdateByIdRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: repo}
}

// Implement
func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req CategoryUpdateReq) error {
	if err := req.validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

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

	if err := hdl.repo.Update(ctx, req.Id, req); err != nil {
		return err
	}

	return nil
}
