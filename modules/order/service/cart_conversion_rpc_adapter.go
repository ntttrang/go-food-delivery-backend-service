package service

import (
	"context"

	"github.com/google/uuid"
	sharerpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// CartConversionRPCAdapter adapts the cart RPC client to the cart conversion service interface
type CartConversionRPCAdapter struct {
	cartRPCClient *sharerpc.CartRPCClient
}

// NewCartToOrderConversionServiceWithRPC creates a new cart conversion service using RPC
func NewCartToOrderConversionServiceWithRPC(cartRPCClient *sharerpc.CartRPCClient) ICartConversionService {
	return &CartConversionRPCAdapter{
		cartRPCClient: cartRPCClient,
	}
}

// ValidateCartForOrder validates that the cart can be used for order creation
func (s *CartConversionRPCAdapter) ValidateCartForOrder(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) error {
	// Use the cart RPC client to validate cart ownership
	err := s.cartRPCClient.ValidateCartOwnership(ctx, cartID, userID)
	if err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}

// ConvertCartToOrderData converts cart items to order data
func (s *CartConversionRPCAdapter) ConvertCartToOrderData(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*OrderCreateDto, error) {
	// Get cart items via RPC
	cartItems, err := s.cartRPCClient.FindCartItemsByCartID(ctx, cartID, userID)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if len(cartItems) == 0 {
		return nil, ordermodel.ErrCartEmpty
	}

	// Validate all items are from the same restaurant
	restaurantID := cartItems[0].RestaurantID
	for _, item := range cartItems {
		if item.RestaurantID != restaurantID {
			return nil, ordermodel.ErrMixedRestaurantItems
		}
	}

	// Convert cart items to order details
	var orderDetails []OrderDetailCreateDto
	var totalPrice float64

	for _, item := range cartItems {
		// For now, we'll use basic food information
		// In a real implementation, you'd fetch food details via RPC
		orderDetail := OrderDetailCreateDto{
			FoodOrigin: &FoodOriginDto{
				Id:          item.FoodID.String(),
				Name:        "Food Item", // Would be fetched from food service
				Description: "",
				Image:       "",
			},
			Price:    0, // Would be fetched from food service
			Quantity: item.Quantity,
			Discount: 0,
		}
		orderDetails = append(orderDetails, orderDetail)
		// totalPrice += orderDetail.Price * float64(orderDetail.Quantity)
	}

	// Create order data
	orderData := &OrderCreateDto{
		UserID:       userID.String(),
		TotalPrice:   totalPrice,
		RestaurantID: restaurantID.String(),
		OrderDetails: orderDetails,
	}

	return orderData, nil
}

// MarkCartAsProcessed marks the cart as processed after successful order creation
func (s *CartConversionRPCAdapter) MarkCartAsProcessed(ctx context.Context, cartID uuid.UUID) error {
	// Use the cart RPC client to update cart status
	err := s.cartRPCClient.UpdateCartStatus(ctx, cartID, "PROCESSED")
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
