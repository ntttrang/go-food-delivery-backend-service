package service

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type CartDetailReq struct {
	UserID uuid.UUID `json:"-"`
	FoodID uuid.UUID `json:"-"`
}

type CartDetailRes struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"userId"`
	FoodID      uuid.UUID           `json:"foodId"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Quantity    int                 `json:"quantity"`
	Status      datatype.CartStatus `json:"status"`
	sharedModel.DateDto
}

// Initialize service
type IGetCartDetailRepo interface {
	FindByUserIdAndFoodId(ctx context.Context, userId, foodId uuid.UUID) (*cartmodel.Cart, error)
}

type IRpcFoodForGetCartRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]cartmodel.Food, error)
}

type GetDetailQueryHandler struct {
	repo        IGetCartDetailRepo
	rpcFoodRepo IRpcFoodForGetCartRepo
}

func NewGetDetailQueryHandler(repo IGetCartDetailRepo, rpcFoodRepo IRpcFoodForGetCartRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{
		repo:        repo,
		rpcFoodRepo: rpcFoodRepo,
	}
}

// Implement
func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req CartDetailReq) (*CartDetailRes, error) {
	if req.UserID == uuid.Nil {
		return nil, datatype.ErrBadRequest.WithWrap(cartmodel.ErrUserIdRequired).WithDebug(cartmodel.ErrUserIdRequired.Error())
	}
	if req.FoodID == uuid.Nil {
		return nil, datatype.ErrBadRequest.WithWrap(cartmodel.ErrFoodIdRequired).WithDebug(cartmodel.ErrFoodIdRequired.Error())
	}

	cart, err := hdl.repo.FindByUserIdAndFoodId(ctx, req.UserID, req.FoodID)
	if err != nil {
		if err == cartmodel.ErrCartNotFound {
			return nil, datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	if cart == nil {
		return nil, datatype.ErrNotFound
	}

	// Check if cart is processed
	if cart.Status == datatype.CartStatusProcessed {
		return nil, datatype.ErrDeleted.WithWrap(cartmodel.ErrCartIsProcessed).WithDebug(cartmodel.ErrCartIsProcessed.Error())
	}

	// Get food detail
	foodMap, err := hdl.rpcFoodRepo.FindByIds(ctx, []uuid.UUID{cart.FoodID})
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return &CartDetailRes{
		ID:          cart.ID,
		UserID:      cart.UserID,
		Name:        foodMap[cart.FoodID].Name,
		Description: foodMap[cart.FoodID].Description,
		Quantity:    cart.Quantity,
		Status:      cart.Status,
		DateDto:     cart.DateDto,
	}, nil
}
