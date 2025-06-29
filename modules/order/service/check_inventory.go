package service

import (
	"context"

	"github.com/google/uuid"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/rpc-client"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"go.opentelemetry.io/otel"
)

// RestaurantAvailability represents restaurant availability info
type RestaurantAvailability struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	IsOpen           bool      `json:"isOpen"`
	AcceptingOrders  bool      `json:"acceptingOrders"`
	ShippingFeePerKm float64   `json:"shippingFeePerKm"`
}

// OrderItem represents an item in the order for inventory checking
type OrderItem struct {
	FoodID   uuid.UUID `json:"foodId"`
	Quantity int       `json:"quantity"`
}

// Repository interfaces
type IInventoryFoodRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]ordermodel.Food, error)
}

type IInventoryRestaurantRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcclient.RPCGetByIdsResponseDTO, error)
}

// Service
type InventoryCheckingService struct {
	foodRepo       IInventoryFoodRepo
	restaurantRepo IInventoryRestaurantRepo
}

func NewInventoryCheckingService(
	foodRepo IInventoryFoodRepo,
	restaurantRepo IInventoryRestaurantRepo,
) *InventoryCheckingService {
	return &InventoryCheckingService{
		foodRepo:       foodRepo,
		restaurantRepo: restaurantRepo,
	}
}

// CheckRestaurantAvailability checks if restaurant is available for orders
func (s *InventoryCheckingService) CheckRestaurantAvailability(ctx context.Context, restaurantID uuid.UUID) error {
	restaurantMap, err := s.restaurantRepo.FindByIds(ctx, []uuid.UUID{restaurantID})
	if err != nil {
		return datatype.ErrNotFound.WithWrap(err).WithDebug("restaurant not found")
	}

	restaurant := restaurantMap[restaurantID]
	// Check if restaurant is active
	if restaurant.Status != "ACTIVE" {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrRestaurantNotAvailable).WithDebug("restaurant is not active")
	}

	return nil
}

// CheckFoodAvailability checks if all food items are available in required quantities
func (s *InventoryCheckingService) CheckFoodAvailability(ctx context.Context, items []OrderItem) error {
	if len(items) == 0 {
		return datatype.ErrBadRequest.WithError("no items to check")
	}

	// Get food IDs
	var foodIDs []uuid.UUID
	for _, item := range items {
		foodIDs = append(foodIDs, item.FoodID)
	}

	// Get food inventory info
	foodMap, err := s.foodRepo.FindByIds(ctx, foodIDs)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug("failed to get food inventory")
	}

	// Check each item
	for _, item := range items {
		food, exists := foodMap[item.FoodID]
		if !exists {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrFoodNotAvailable).WithDebug("food not found: " + item.FoodID.String())
		}

		// Check if food is active
		if food.Status != "ACTIVE" {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrFoodNotAvailable).WithDebug("food is not active: " + food.Name)
		}
	}

	return nil
}

// CheckOrderInventory performs comprehensive inventory check for an order
func (s *InventoryCheckingService) CheckOrderInventory(ctx context.Context, restaurantID uuid.UUID, items []OrderItem) error {
	_, dbSpanCkInvtry := otel.Tracer("").Start(ctx, "check-inventory")
	defer dbSpanCkInvtry.End()

	// Check restaurant availability: check status = ACTIVE
	if err := s.CheckRestaurantAvailability(ctx, restaurantID); err != nil {
		return err
	}

	// Check food availability: check food exist and status = ACTIVE
	if err := s.CheckFoodAvailability(ctx, items); err != nil {
		return err
	}

	// Validate all items belong to the same restaurant
	if err := s.validateItemsRestaurant(ctx, restaurantID, items); err != nil {
		return err
	}

	return nil
}

// validateItemsRestaurant ensures all food items belong to the specified restaurant
func (s *InventoryCheckingService) validateItemsRestaurant(ctx context.Context, restaurantID uuid.UUID, items []OrderItem) error {
	var foodIDs []uuid.UUID
	for _, item := range items {
		foodIDs = append(foodIDs, item.FoodID)
	}

	foodMap, err := s.foodRepo.FindByIds(ctx, foodIDs)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug("failed to validate food restaurant")
	}

	for _, item := range items {
		food, exists := foodMap[item.FoodID]
		if !exists {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrFoodNotAvailable).WithDebug("food not found: " + item.FoodID.String())
		}

		if food.RestaurantId != restaurantID {
			return datatype.ErrBadRequest.WithError("food item does not belong to the specified restaurant: " + food.Name)
		}
	}

	return nil
}
