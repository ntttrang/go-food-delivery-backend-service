package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type CategoryDeleteReq struct {
	Id        uuid.UUID
	Requester datatype.Requester `json:"-"`
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
	if req.Id == uuid.Nil {
		return datatype.ErrBadRequest.WithDebug(categorymodel.ErrIdRequired.Error())
	}
	// Authorization check
	if req.Requester == nil {
		return datatype.ErrUnauthorized.WithDebug(categorymodel.ErrRequesterRequired.Error())
	}

	// Check if requester is admin or user (not shipper)
	// Use type assertion to get the role as a string
	role := req.Requester.GetRole()

	// Check role
	isAuthorized := role == string(datatype.RoleAdmin) || role == string(datatype.RoleUser)
	if !isAuthorized {
		return datatype.ErrForbidden.WithDebug(categorymodel.ErrPermission.Error())
	}

	category, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, categorymodel.ErrCategoryNotFound) {
			return datatype.ErrNotFound.WithDebug(categorymodel.ErrCategoryNotFound.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if category.Status == string(datatype.StatusDeleted) {
		return datatype.ErrDeleted.WithError(categorymodel.ErrCategoryIsDeleted.Error())
	}

	if err := hdl.repo.Delete(ctx, req.Id); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil

}
