package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IListRestaurantCommentsRepo interface {
	FindByRestaurantId(ctx context.Context, restaurantId string) ([]restaurantmodel.RestaurantRating, error)
}

type IRestaurantRepo interface {
	FindRestaurantByIds(ctx context.Context, ids []uuid.UUID) ([]restaurantmodel.Restaurant, error)
}

type ListRestaurantCommentsQueryHandler struct {
	repo       IListRestaurantCommentsRepo
	reststRepo IRestaurantRepo
}

func NewListRestaurantCommentsQueryHandler(repo IListRestaurantCommentsRepo, reststRepo IRestaurantRepo) *ListRestaurantCommentsQueryHandler {
	return &ListRestaurantCommentsQueryHandler{
		repo:       repo,
		reststRepo: reststRepo,
	}
}

func (hdl *ListRestaurantCommentsQueryHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantRatingListReq) ([]restaurantmodel.RestaurantRatingListRes, error) {
	restaurantRatings, err := hdl.repo.FindByRestaurantId(ctx, req.RestaurantId)
	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}

		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var restaurantIds []uuid.UUID
	var userIds []uuid.UUID
	for _, rr := range restaurantRatings {
		restaurantIds = append(restaurantIds, rr.RestaurantID)
		userIds = append(userIds, rr.UserID)
	}

	restaurants, err := hdl.reststRepo.FindRestaurantByIds(ctx, restaurantIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	restaurantMap := make(map[uuid.UUID]string, len(restaurants))
	for _, r := range restaurants {
		restaurantMap[r.Id] = r.Name
	}

	var resp []restaurantmodel.RestaurantRatingListRes
	for _, rr := range restaurantRatings {
		var rs restaurantmodel.RestaurantRatingListRes
		rs.RestaurantId = rr.RestaurantID
		rs.RestaurantName = restaurantMap[rr.RestaurantID]
		resp = append(resp, rs)
	}
	return resp, nil
}
