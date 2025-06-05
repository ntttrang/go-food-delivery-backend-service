package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type CartUpdateReq struct {
	ID     uuid.UUID `json:"-"`
	Status *string   `json:"status,omitempty"`
}

func (c *CartUpdateReq) Validate() error {
	if c.ID == uuid.Nil {
		return cartmodel.ErrCartIdRequired
	}

	if c.Status != nil {
		status := *c.Status
		if status != datatype.CartStatusActive && status != datatype.CartStatusUpdated && status != datatype.CartStatusProcessed {
			return cartmodel.ErrCartStatusInvalid
		}
	}

	return nil
}

// Initialize service
type IUpdateCartRepo interface {
	FindById(ctx context.Context, id uuid.UUID) ([]cartmodel.Cart, error)
	UpdateStatusById(ctx context.Context, id uuid.UUID, status string) error
}

type UpdateCommandHandler struct {
	repo IUpdateCartRepo
}

func NewUpdateCommandHandler(repo IUpdateCartRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: repo}
}

// Implement
func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req CartUpdateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Check if cart exists
	carts, err := hdl.repo.FindById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, cartmodel.ErrCartNotFound) {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	if len(carts) == 0 {
		return datatype.ErrNotFound
	}

	// Check if cart is already processed
	if carts[0].Status == datatype.CartStatusProcessed {
		return datatype.ErrDeleted.WithWrap(cartmodel.ErrCartIsProcessed).WithDebug(cartmodel.ErrCartIsProcessed.Error())
	}

	// Update cart in database
	if err := hdl.repo.UpdateStatusById(ctx, req.ID, *req.Status); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
