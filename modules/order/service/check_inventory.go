package service

import (
	"context"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// FoodInventory represents food item with availability info
type FoodInventory struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurantId"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"`
	Available    bool      `json:"available"`
}

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
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]FoodInventory, error)
	CheckAvailability(ctx context.Context, foodID uuid.UUID, quantity int) (bool, error)
}

type IInventoryRestaurantRepo interface {
	FindByID(ctx context.Context, id uuid.UUID) (*RestaurantAvailability, error)
	IsAcceptingOrders(ctx context.Context, restaurantID uuid.UUID) (bool, error)
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
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return datatype.ErrNotFound.WithWrap(err).WithDebug("restaurant not found")
	}

	// Check if restaurant is active
	if restaurant.Status != "ACTIVE" {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrRestaurantNotAvailable).WithDebug("restaurant is not active")
	}

	// Check if restaurant is accepting orders
	acceptingOrders, err := s.restaurantRepo.IsAcceptingOrders(ctx, restaurantID)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug("failed to check restaurant order acceptance")
	}

	if !acceptingOrders {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrRestaurantNotAvailable).WithDebug("restaurant is not accepting orders")
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

		// Check if food is available
		if !food.Available {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrFoodNotAvailable).WithDebug("food is not available: " + food.Name)
		}

		// Check specific quantity availability
		available, err := s.foodRepo.CheckAvailability(ctx, item.FoodID, item.Quantity)
		if err != nil {
			return datatype.ErrInternalServerError.WithWrap(err).WithDebug("failed to check food quantity availability")
		}

		if !available {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrInventoryInsufficient).WithDebug("insufficient quantity for: " + food.Name)
		}
	}

	return nil
}

// CheckOrderInventory performs comprehensive inventory check for an order
func (s *InventoryCheckingService) CheckOrderInventory(ctx context.Context, restaurantID uuid.UUID, items []OrderItem) error {
	// Check restaurant availability first
	if err := s.CheckRestaurantAvailability(ctx, restaurantID); err != nil {
		return err
	}

	// Check food availability
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

		if food.RestaurantID != restaurantID {
			return datatype.ErrBadRequest.WithError("food item does not belong to the specified restaurant: " + food.Name)
		}
	}

	return nil
}

// GetRestaurantInfo gets restaurant information for order processing
func (s *InventoryCheckingService) GetRestaurantInfo(ctx context.Context, restaurantID uuid.UUID) (*RestaurantAvailability, error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, datatype.ErrNotFound.WithWrap(err).WithDebug("restaurant not found")
	}

	return restaurant, nil
}
