package service

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type CartUpsertDto struct {
	FoodID   uuid.UUID `json:"foodId"`
	Quantity int       `json:"quantity"`

	UserID uuid.UUID `json:"-"` // Get from token
	ID     uuid.UUID `json:"-"` // Internal BE
}

func (c *CartUpsertDto) Validate() error {
	if c.UserID == uuid.Nil {
		return cartmodel.ErrUserIdRequired
	}

	if c.FoodID == uuid.Nil {
		return cartmodel.ErrFoodIdRequired
	}

	if c.Quantity <= 0 {
		return cartmodel.ErrQuantityInvalid
	}

	return nil
}

func (c CartUpsertDto) ConvertToCart() *cartmodel.Cart {
	return &cartmodel.Cart{
		UserID:   c.UserID,
		FoodID:   c.FoodID,
		Quantity: c.Quantity,
	}
}

// Initialize service
type ICreateCartRepository interface {
	Insert(ctx context.Context, cart *cartmodel.Cart) error
	FindByUserIdAndFoodId(ctx context.Context, userId, foodId uuid.UUID) (*cartmodel.Cart, error)
	FindByUserIdAndRestaurantId(ctx context.Context, userId, foodId uuid.UUID) ([]cartmodel.Cart, error)
	Update(ctx context.Context, cart *cartmodel.Cart) error
}

type IFoodRepository interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]cartmodel.Food, error)
}

type CreateCommandHandler struct {
	repo     ICreateCartRepository
	foodRepo IFoodRepository
}

func NewCreateCommandHandler(repo ICreateCartRepository, foodRepo IFoodRepository) *CreateCommandHandler {
	return &CreateCommandHandler{repo: repo, foodRepo: foodRepo}
}

// Implement
func (s *CreateCommandHandler) Execute(ctx context.Context, data *CartUpsertDto) error {
	if err := data.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Check food
	var ids []uuid.UUID
	ids = append(ids, data.FoodID)
	food, err := s.foodRepo.FindByIds(ctx, ids)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	if food != nil && food[data.FoodID].RestaurantId == uuid.Nil {
		return datatype.ErrBadRequest.WithWrap(cartmodel.ErrFoodInvalid).WithDebug(cartmodel.ErrFoodInvalid.Error())
	}

	// Check if cart already exists for this user and food
	existingCart, err := s.repo.FindByUserIdAndFoodId(ctx, data.UserID, data.FoodID)
	if err == nil && existingCart != nil {
		// Cart exists, update quantity instead
		existingCart.Quantity += data.Quantity
		existingCart.Status = cartmodel.CartStatusUpdated
		if err := s.repo.Update(ctx, existingCart); err != nil {
			return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
		data.ID = existingCart.ID
		return nil
	}

	// Create new cart
	cart := data.ConvertToCart()

	// Check if cart already exists for this user and restaurant
	var cartId uuid.UUID
	existingCartByRestaurants, err := s.repo.FindByUserIdAndRestaurantId(ctx, data.UserID, food[data.FoodID].RestaurantId)
	if err == nil && len(existingCartByRestaurants) > 0 {
		cartId = existingCartByRestaurants[0].ID
	} else {
		cartId, _ = uuid.NewV7()
	}
	cart.ID = cartId

	cart.Status = cartmodel.CartStatusActive // Always set Active Status when insert
	cart.RestaurantId = food[data.FoodID].RestaurantId

	if err := s.repo.Insert(ctx, cart); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	data.ID = cart.ID

	return nil
}
