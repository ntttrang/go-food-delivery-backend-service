package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListRestaurantCommentsRepo interface {
	FindByRestaurantIdOrUserId(ctx context.Context, req restaurantmodel.RestaurantRatingListReq) ([]restaurantmodel.RestaurantRating, int64, error)
}

type IUserRPCClientRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.User, error)
}

type ListRestaurantCommentsQueryHandler struct {
	restrepo IListRestaurantCommentsRepo
	userRepo IUserRPCClientRepo
}

func NewListRestaurantCommentsQueryHandler(restrepo IListRestaurantCommentsRepo, userRepo IUserRPCClientRepo) *ListRestaurantCommentsQueryHandler {
	return &ListRestaurantCommentsQueryHandler{
		restrepo: restrepo,
		userRepo: userRepo,
	}
}

func (hdl *ListRestaurantCommentsQueryHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantRatingListReq) (*restaurantmodel.RestaurantRatingListRes, error) {
	restaurantRatings, total, err := hdl.restrepo.FindByRestaurantIdOrUserId(ctx, req)
	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}

		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var userIds []uuid.UUID
	for _, rr := range restaurantRatings {
		userIds = append(userIds, rr.UserID)
	}

	userMap, err := hdl.userRepo.FindByIds(ctx, userIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var rsratings []restaurantmodel.RestaurantRatingListDto
	for _, rr := range restaurantRatings {
		var rs restaurantmodel.RestaurantRatingListDto
		rs.Id = rr.ID
		rs.UserId = rr.UserID
		rs.FirstName = userMap[rr.UserID].FirstName
		rs.LastName = userMap[rr.UserID].LastName
		rs.RestaurantId = rr.RestaurantID
		rs.Comment = rr.Comment
		rs.Rating = rr.Point
		rs.CreatedAt = rr.CreatedAt
		rsratings = append(rsratings, rs)
	}

	var resp restaurantmodel.RestaurantRatingListRes
	resp.Items = rsratings
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
