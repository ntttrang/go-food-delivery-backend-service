package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type FoodRatingListReq struct {
	FoodId string `json:"foodId" form:"foodId"`
	UserId string `json:"userId" form:"userId"`
	sharedModel.PagingDto
}

func (r *FoodRatingListReq) Validate() error {
	if r.FoodId == "" && r.UserId == "" {
		return foodmodel.ErrFieldRequired
	}

	return nil
}

type FoodRatingListDto struct {
	Id        uuid.UUID  `json:"id"`
	FoodId    uuid.UUID  `json:"restaurantId"`
	UserId    uuid.UUID  `json:"userId"`
	FirstName string     `json:"frstName"`
	LastName  string     `json:"lastName"`
	Avatar    *string    `json:"avatar"`
	Rating    float64    `json:"rating"`
	Comment   *string    `json:"comment"`
	CreatedAt *time.Time `json:"createdAt"`
}

type FoodRatingListRes struct {
	Items      []FoodRatingListDto   `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

// Initilize service
type IListFoodCommentsRepo interface {
	FindByFoodIdOrUserId(ctx context.Context, req FoodRatingListReq) ([]foodmodel.FoodRatings, int64, error)
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

// Implement
func (hdl *ListFoodCommentsQueryHandler) Execute(ctx context.Context, req FoodRatingListReq) (*FoodRatingListRes, error) {
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

	var rsratings []FoodRatingListDto
	for _, fr := range foodRatings {
		var rs FoodRatingListDto
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

	var resp FoodRatingListRes
	resp.Items = rsratings
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
