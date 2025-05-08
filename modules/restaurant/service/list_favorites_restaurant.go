package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FavoriteRestaurantListReq struct {
	UserId uuid.UUID `json:"-" form:"-"`
	sharedModel.PagingDto
}

type IListFavoritesRestaurantRepo interface {
	FindFavRestaurant(ctx context.Context, req FavoriteRestaurantListReq) ([]RestaurantSearchResDto, int64, error)
}

type ListFavoritesRestaurantQueryHandler struct {
	repo IListFavoritesRestaurantRepo
}

func NewGetFavoritesRestaurantQueryHandler(repo IListFavoritesRestaurantRepo) *ListFavoritesRestaurantQueryHandler {
	return &ListFavoritesRestaurantQueryHandler{repo: repo}
}

func (hdl *ListFavoritesRestaurantQueryHandler) Execute(ctx context.Context, req FavoriteRestaurantListReq) (*RestaurantSearchRes, error) {
	restaurants, total, err := hdl.repo.FindFavRestaurant(ctx, req)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp RestaurantSearchRes
	resp.Items = restaurants
	if restaurants == nil {
		resp.Items = make([]RestaurantSearchResDto, 0)
	}

	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
