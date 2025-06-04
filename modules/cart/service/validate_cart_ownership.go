package service

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs for cart ownership validation
type CartValidationReq struct {
	CartID uuid.UUID `json:"cartId"`
	UserID uuid.UUID `json:"userId"`
}

type CartValidationRes struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

// Repository interface for cart validation
type IValidateCartOwnershipRepo interface {
	FindCartByCartIDAndUserID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*cartmodel.Cart, error)
}

// Service handler
type ValidateCartOwnershipQueryHandler struct {
	repo IValidateCartOwnershipRepo
}

func NewValidateCartOwnershipQueryHandler(repo IValidateCartOwnershipRepo) *ValidateCartOwnershipQueryHandler {
	return &ValidateCartOwnershipQueryHandler{
		repo: repo,
	}
}

// Execute validates cart ownership
func (hdl *ValidateCartOwnershipQueryHandler) Execute(ctx context.Context, req CartValidationReq) (*CartValidationRes, error) {
	// Validate input
	if req.CartID == uuid.Nil {
		return nil, datatype.ErrBadRequest.WithError("cart ID is required")
	}

	if req.UserID == uuid.Nil {
		return nil, datatype.ErrBadRequest.WithError("user ID is required")
	}

	// Check if cart exists and belongs to user
	cart, err := hdl.repo.FindCartByCartIDAndUserID(ctx, req.CartID, req.UserID)
	if err != nil {
		if err == cartmodel.ErrCartNotFound {
			return &CartValidationRes{
				Valid:  false,
				Reason: "cart not found or does not belong to user",
			}, nil
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Check if cart is already processed
	if cart.Status == datatype.CartStatusProcessed {
		return &CartValidationRes{
			Valid:  false,
			Reason: "cart has already been processed",
		}, nil
	}

	// Cart is valid
	return &CartValidationRes{
		Valid: true,
	}, nil
}
