package service

import (
	"context"

	"github.com/google/uuid"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/rpc-client"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Repository interfaces
type ICartConversionRepo interface {
	FindById(ctx context.Context, cartId string, userId string) ([]rpcclient.CartSummaryData, error)
	UpdateCartStatus(ctx context.Context, cartID uuid.UUID, status string) error
}

type IFoodConversionRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]ordermodel.Food, error)
}

type IRestaurantConversionRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcclient.RPCGetByIdsResponseDTO, error)
}

// Cart conversion service interface
type ICartConversionService interface {
	ValidateCartForOrder(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) error
	ConvertCartToOrderData(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*OrderCreateDto, error)
	MarkCartAsProcessed(ctx context.Context, cartID uuid.UUID) error
}

// Service
type CartToOrderConversionService struct {
	cartRepo       ICartConversionRepo
	foodRepo       IFoodConversionRepo
	restaurantRepo IRestaurantConversionRepo
}

func NewCartToOrderConversionService(
	cartRepo ICartConversionRepo,
	foodRepo IFoodConversionRepo,
	restaurantRepo IRestaurantConversionRepo,
) *CartToOrderConversionService {
	return &CartToOrderConversionService{
		cartRepo:       cartRepo,
		foodRepo:       foodRepo,
		restaurantRepo: restaurantRepo,
	}
}

// ConvertCartToOrderData converts cart items to order data
func (s *CartToOrderConversionService) ConvertCartToOrderData(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*OrderCreateDto, error) {
	// Get cart items
	cartItems, err := s.cartRepo.FindById(ctx, cartID.String(), userID.String())
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if len(cartItems) == 0 {
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrCartEmpty).WithDebug(ordermodel.ErrCartEmpty.Error())
	}

	// Check if cart is already processed
	for _, item := range cartItems {
		if item.Status == datatype.CartStatusProcessed {
			return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrCartAlreadyProcessed).WithDebug(ordermodel.ErrCartAlreadyProcessed.Error())
		}
	}

	// Validate all items are from the same restaurant
	restaurantID := cartItems[0].RestaurantID
	for _, item := range cartItems {
		if item.RestaurantID != restaurantID {
			return nil, datatype.ErrBadRequest.WithError("all cart items must be from the same restaurant")
		}
	}

	// Get food details
	var foodIDs []uuid.UUID
	for _, item := range cartItems {
		foodIDs = append(foodIDs, item.FoodId)
	}

	foodMap, err := s.foodRepo.FindByIds(ctx, foodIDs)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Get restaurant details
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantID})
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	restaurant := restaurantMap[restaurantID]
	// Check restaurant availability
	if restaurant.Status != "ACTIVE" {
		return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrRestaurantNotAvailable).WithDebug(ordermodel.ErrRestaurantNotAvailable.Error())
	}

	// Convert cart items to order details
	var orderDetails []OrderDetailCreateDto
	var totalPrice float64

	for _, item := range cartItems {
		food, exists := foodMap[item.FoodId]
		if !exists {
			return nil, datatype.ErrBadRequest.WithWrap(ordermodel.ErrFoodNotAvailable).WithDebug("food not found: " + item.FoodId.String())
		}

		// Create food origin JSON
		foodOrigin := &FoodOriginDto{
			Id:          food.Id.String(),
			Name:        food.Name,
			Description: food.Description,
			Image:       food.Images,
		}

		orderDetail := OrderDetailCreateDto{
			FoodOrigin: foodOrigin,
			Price:      food.Price,
			Quantity:   item.Quantity,
			Discount:   0, // Default no discount
		}

		orderDetails = append(orderDetails, orderDetail)
		totalPrice += food.Price * float64(item.Quantity)
	}

	// Create order DTO
	orderDto := &OrderCreateDto{
		UserID:       userID.String(),
		TotalPrice:   totalPrice,
		RestaurantID: restaurantID.String(),
		OrderDetails: orderDetails,
	}

	return orderDto, nil
}

// MarkCartAsProcessed marks the cart as processed after successful order creation
func (s *CartToOrderConversionService) MarkCartAsProcessed(ctx context.Context, cartID uuid.UUID) error {
	err := s.cartRepo.UpdateCartStatus(ctx, cartID, datatype.CartStatusProcessed)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}

// ValidateCartForOrder validates that cart can be converted to order
func (s *CartToOrderConversionService) ValidateCartForOrder(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) error {
	cartItems, err := s.cartRepo.FindById(ctx, cartID.String(), userID.String())
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if len(cartItems) == 0 {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrCartEmpty).WithDebug(ordermodel.ErrCartEmpty.Error())
	}

	// Check if cart is already processed
	for _, item := range cartItems {
		if item.Status == datatype.CartStatusProcessed {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrCartAlreadyProcessed).WithDebug(ordermodel.ErrCartAlreadyProcessed.Error())
		}
	}

	// Validate all items are from the same restaurant
	restaurantID := cartItems[0].RestaurantID
	for _, item := range cartItems {
		if item.RestaurantID != restaurantID {
			return datatype.ErrBadRequest.WithError("all cart items must be from the same restaurant")
		}
	}

	return nil
}
