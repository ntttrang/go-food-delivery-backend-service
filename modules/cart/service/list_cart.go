package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/rpcclient"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	sharedComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type CartSearchDto struct {
	UserID *string `json:"userId" form:"userId"`
}

type CartListReq struct {
	CartSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

// No need to paginate
// Maximum 10 carts
type CartListRes struct {
	Items []CartDto `json:"items"`
}

type CartDto struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"userId"`
	RestaurantId   uuid.UUID `json:"restaurantId"`
	RestaurantName string    `json:"restaurantName"`
	Quantity       int       `json:"quantity"`
	EstDistance    float64   `json:"estDistance"` // Km
	EstTime        float64   `json:"estTime"`     // minute

	DropOffLat float64 `json:"-"`
	DropOffLng float64 `json:"-"`
}

// Initialize service
type IListCartRepo interface {
	ListByUserId(ctx context.Context, req CartListReq) ([]cartmodel.Cart, error)
}

type IRpcRestaurantRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcclient.RPCGetByIdsResponseDTO, error)
}

type ListQueryHandler struct {
	cartRepo          IListCartRepo
	rpcRestaurantRepo IRpcRestaurantRepo
}

func NewListQueryHandler(cartRepo IListCartRepo, rpcRestaurantRepo IRpcRestaurantRepo) *ListQueryHandler {
	return &ListQueryHandler{
		cartRepo:          cartRepo,
		rpcRestaurantRepo: rpcRestaurantRepo,
	}
}

// Implement
func (hdl *ListQueryHandler) Execute(ctx context.Context, req CartListReq) (*CartListRes, error) {
	carts, err := hdl.cartRepo.ListByUserId(ctx, req)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var items []CartDto
	var restaurantIds []uuid.UUID
	for _, cart := range carts {
		items = append(items, CartDto{
			ID:           cart.ID,
			UserID:       cart.UserID,
			RestaurantId: cart.RestaurantId,
			Quantity:     cart.Quantity,
		})
		restaurantIds = append(restaurantIds, cart.RestaurantId)
	}

	restaurantMap, err := hdl.rpcRestaurantRepo.FindByIds(ctx, restaurantIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	for i := 0; i < len(items); i++ {
		items[i].RestaurantName = restaurantMap[items[i].RestaurantId].Name

		// Estimate distance & time
		lat := restaurantMap[items[i].RestaurantId].Lat
		lng := restaurantMap[items[i].RestaurantId].Lng
		currentLat := carts[i].DropOffLat
		currentLng := carts[i].DropOffLng

		distance := sharedComponent.Haversine(currentLat, currentLng, lat, lng)

		items[i].EstDistance = distance
	}

	var resp CartListRes
	resp.Items = items
	return &resp, nil
}
