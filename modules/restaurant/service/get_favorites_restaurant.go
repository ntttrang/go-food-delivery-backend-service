package service

import "context"

type IGetFavoritesRestaurantRepo interface {
}

type GetFavoritesRestaurantQueryHandler struct {
	repo IGetFavoritesRestaurantRepo
}

func NewGetFavoritesRestaurantQueryHandler(repo IGetFavoritesRestaurantRepo) *GetFavoritesRestaurantQueryHandler {
	return &GetFavoritesRestaurantQueryHandler{repo: repo}
}

func (hdl *GetFavoritesRestaurantQueryHandler) Execute(ctx context.Context, userId string) {

}
