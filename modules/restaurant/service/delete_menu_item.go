package service

import (
	"context"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IDeleteMenuItemRepo interface {
	Delete(ctx context.Context, req *MenuItemCreateReq) error
}

type DeleteMenuItemCommandHandler struct {
	menuItemRepo IDeleteMenuItemRepo
}

func NewDeleteMenuItemCommandHandler(menuItemRepo IDeleteMenuItemRepo) *DeleteMenuItemCommandHandler {
	return &DeleteMenuItemCommandHandler{
		menuItemRepo: menuItemRepo,
	}
}

func (s *DeleteMenuItemCommandHandler) Execute(ctx context.Context, req *MenuItemCreateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	if err := s.menuItemRepo.Delete(ctx, req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
