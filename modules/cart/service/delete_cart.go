package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type CartDeleteReq struct {
	UserID uuid.UUID `json:"-"`
	FoodID uuid.UUID `json:"-"`
}

// Initialize service
type IDeleteCartRepo interface {
	FindByUserIdAndFoodId(ctx context.Context, userId, foodId uuid.UUID) (*cartmodel.Cart, error)
	Delete(ctx context.Context, userId, foodId uuid.UUID) error
}

type DeleteCommandHandler struct {
	repo IDeleteCartRepo
}

func NewDeleteCommandHandler(repo IDeleteCartRepo) *DeleteCommandHandler {
	return &DeleteCommandHandler{repo: repo}
}

// Implement
func (hdl *DeleteCommandHandler) Execute(ctx context.Context, req CartDeleteReq) error {
	if req.UserID == uuid.Nil {
		return datatype.ErrBadRequest.WithWrap(cartmodel.ErrUserIdRequired).WithDebug(cartmodel.ErrUserIdRequired.Error())
	}
	if req.FoodID == uuid.Nil {
		return datatype.ErrBadRequest.WithWrap(cartmodel.ErrFoodIdRequired).WithDebug(cartmodel.ErrFoodIdRequired.Error())
	}

	// Check if cart exists
	cart, err := hdl.repo.FindByUserIdAndFoodId(ctx, req.UserID, req.FoodID)
	if err != nil {
		if errors.Is(err, cartmodel.ErrCartNotFound) {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Check if cart is processed
	if cart.Status == datatype.CartStatusProcessed {
		return datatype.ErrNotFound.WithWrap(cartmodel.ErrCartIsProcessed).WithDebug(cartmodel.ErrCartIsProcessed.Error())
	}

	// Delete cart (soft delete by updating status)
	if err := hdl.repo.Delete(ctx, req.UserID, req.FoodID); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
