package service

// import (
// 	"context"

// 	"github.com/google/uuid"
// 	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
// 	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
// )

// // Define DTOs for cart summary
// type CartSummaryReq struct {
// 	CartID uuid.UUID `json:"cartId"`
// 	UserID uuid.UUID `json:"userId"`
// }

// type CartSummaryRes struct {
// 	CartID       uuid.UUID `json:"cartId"`
// 	UserID       uuid.UUID `json:"userId"`
// 	RestaurantID uuid.UUID `json:"restaurantId"`
// 	ItemCount    int       `json:"itemCount"`
// 	TotalPrice   float64   `json:"totalPrice"`
// 	Status       string    `json:"status"`
// 	CreatedAt    string    `json:"createdAt"`
// 	UpdatedAt    string    `json:"updatedAt"`
// }

// // Repository interface for cart summary
// type IGetCartSummaryRepo interface {
// 	GetCartSummaryByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*CartSummaryData, error)
// }

// // Data structure for cart summary from repository
// type CartSummaryData struct {
// 	CartID       uuid.UUID
// 	UserID       uuid.UUID
// 	RestaurantID uuid.UUID
// 	ItemCount    int
// 	TotalPrice   float64
// 	Status       string
// 	CreatedAt    string
// 	UpdatedAt    string
// }

// // Service handler
// type GetCartSummaryQueryHandler struct {
// 	repo IGetCartSummaryRepo
// }

// func NewGetCartSummaryQueryHandler(repo IGetCartSummaryRepo) *GetCartSummaryQueryHandler {
// 	return &GetCartSummaryQueryHandler{
// 		repo: repo,
// 	}
// }

// // Execute gets cart summary
// func (hdl *GetCartSummaryQueryHandler) Execute(ctx context.Context, req CartSummaryReq) (*CartSummaryRes, error) {
// 	// Validate input
// 	if req.CartID == uuid.Nil {
// 		return nil, datatype.ErrBadRequest.WithError("cart ID is required")
// 	}

// 	if req.UserID == uuid.Nil {
// 		return nil, datatype.ErrBadRequest.WithError("user ID is required")
// 	}

// 	// Get cart summary from repository
// 	summary, err := hdl.repo.GetCartSummaryByCartID(ctx, req.CartID, req.UserID)
// 	if err != nil {
// 		if err == cartmodel.ErrCartNotFound {
// 			return nil, datatype.ErrNotFound.WithError("cart not found")
// 		}
// 		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
// 	}

// 	// Convert to response DTO
// 	return &CartSummaryRes{
// 		CartID:       summary.CartID,
// 		UserID:       summary.UserID,
// 		RestaurantID: summary.RestaurantID,
// 		ItemCount:    summary.ItemCount,
// 		TotalPrice:   summary.TotalPrice,
// 		Status:       summary.Status,
// 		CreatedAt:    summary.CreatedAt,
// 		UpdatedAt:    summary.UpdatedAt,
// 	}, nil
// }
