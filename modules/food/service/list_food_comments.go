package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListFoodCommentsRepo interface {
	FindByFoodIdOrUserId(ctx context.Context, req foodmodel.FoodRatingListReq) ([]foodmodel.FoodRatings, int64, error)
}

type IUserRPCClientRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]foodmodel.User, error)
}

type ListFoodCommentsQueryHandler struct {
	restrepo IListFoodCommentsRepo
	userRepo IUserRPCClientRepo
}

func NewListFoodCommentsQueryHandler(restrepo IListFoodCommentsRepo, userRepo IUserRPCClientRepo) *ListFoodCommentsQueryHandler {
	return &ListFoodCommentsQueryHandler{
		restrepo: restrepo,
		userRepo: userRepo,
	}
}

func (hdl *ListFoodCommentsQueryHandler) Execute(ctx context.Context, req foodmodel.FoodRatingListReq) (*foodmodel.FoodRatingListRes, error) {
	foodRatings, total, err := hdl.restrepo.FindByFoodIdOrUserId(ctx, req)
	if err != nil {
		if errors.Is(err, foodmodel.ErrFoodNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(foodmodel.ErrFoodNotFound.Error())
		}

		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var userIds []uuid.UUID
	for _, rr := range foodRatings {
		userIds = append(userIds, rr.UserId)
	}

	userMap, err := hdl.userRepo.FindByIds(ctx, userIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var rsratings []foodmodel.FoodRatingListDto
	for _, fr := range foodRatings {
		var rs foodmodel.FoodRatingListDto
		rs.Id = fr.Id
		rs.UserId = fr.UserId
		rs.FirstName = userMap[fr.UserId].FirstName
		rs.LastName = userMap[fr.UserId].LastName
		rs.FoodId = fr.FoodId
		rs.Comment = &fr.Comment
		rs.Rating = fr.Point
		rs.CreatedAt = fr.CreatedAt
		rsratings = append(rsratings, rs)
	}

	var resp foodmodel.FoodRatingListRes
	resp.Items = rsratings
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
