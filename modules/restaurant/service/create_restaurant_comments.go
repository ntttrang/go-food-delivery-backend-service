package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

type IInsertCommentRestaurantRepo interface {
	Insert(ctx context.Context, req *restaurantmodel.RestaurantCommentCreateReq) error
}

type CreateRestaurantCommentCommandHandler struct {
	repo IInsertCommentRestaurantRepo
}

func NewCommentRestaurantCommandHandler(repo IInsertCommentRestaurantRepo) *CreateRestaurantCommentCommandHandler {
	return &CreateRestaurantCommentCommandHandler{repo: repo}
}

func (hdl *CreateRestaurantCommentCommandHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantCommentCreateReq) error {
	if err := req.Validate(); err != nil {
		return err
	}

	req.Id, _ = uuid.NewV7()
	if err := hdl.repo.Insert(ctx, &req); err != nil {
		return err
	}
	return nil
}
