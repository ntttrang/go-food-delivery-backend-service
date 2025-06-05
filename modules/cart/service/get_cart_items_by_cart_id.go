package service

// import (
// 	"context"

// 	"github.com/google/uuid"
// 	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
// 	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
// )

// // Define DTOs for cart items by cart ID
// type CartItemsByCartIDReq struct {
// 	CartID uuid.UUID `json:"cartId"`
// 	UserID uuid.UUID `json:"userId"`
// }

// type CartItemsByCartIDRes struct {
// 	Items []CartItemByCartIDDto `json:"items"`
// }

// type CartItemByCartIDDto struct {
// 	ID           uuid.UUID `json:"id"`
// 	UserID       uuid.UUID `json:"userId"`
// 	FoodID       uuid.UUID `json:"foodId"`
// 	RestaurantID uuid.UUID `json:"restaurantId"`
// 	Quantity     int       `json:"quantity"`
// 	Status       string    `json:"status"`
// 	CreatedAt    string    `json:"createdAt"`
// 	UpdatedAt    string    `json:"updatedAt"`
// }

// // Repository interface for cart items by cart ID
// type IGetCartItemsByCartIDRepo interface {
// 	FindCartItemsByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) ([]cartmodel.Cart, error)
// }

// // Service handler
// type GetCartItemsByCartIDQueryHandler struct {
// 	repo IGetCartItemsByCartIDRepo
// }

// func NewGetCartItemsByCartIDQueryHandler(repo IGetCartItemsByCartIDRepo) *GetCartItemsByCartIDQueryHandler {
// 	return &GetCartItemsByCartIDQueryHandler{
// 		repo: repo,
// 	}
// }

// // Execute retrieves cart items by cart ID
// func (hdl *GetCartItemsByCartIDQueryHandler) Execute(ctx context.Context, req CartItemsByCartIDReq) (*CartItemsByCartIDRes, error) {
// 	// Validate input
// 	if req.CartID == uuid.Nil {
// 		return nil, datatype.ErrBadRequest.WithError("cart ID is required")
// 	}

// 	if req.UserID == uuid.Nil {
// 		return nil, datatype.ErrBadRequest.WithError("user ID is required")
// 	}

// 	// Get cart items from repository
// 	carts, err := hdl.repo.FindCartItemsByCartID(ctx, req.CartID, req.UserID)
// 	if err != nil {
// 		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
// 	}

// 	// Convert to response DTOs
// 	var items []CartItemByCartIDDto
// 	for _, cart := range carts {
// 		items = append(items, CartItemByCartIDDto{
// 			ID:           cart.ID,
// 			UserID:       cart.UserID,
// 			FoodID:       cart.FoodID,
// 			RestaurantID: cart.RestaurantId,
// 			Quantity:     cart.Quantity,
// 			Status:       cart.Status,
// 			CreatedAt:    cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 			UpdatedAt:    cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 		})
// 	}

// 	return &CartItemsByCartIDRes{
// 		Items: items,
// 	}, nil
// }
