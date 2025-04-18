package service

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListRep interface {
	ListFoods(ctx context.Context, req foodmodel.ListFoodReq) ([]foodmodel.Food, int64, error)
}

type ListCommandHandler struct {
	repo IListRep
}

func NewListCommandHandler(repo IListRep) *ListCommandHandler {
	return &ListCommandHandler{
		repo: repo,
	}
}

func (s *ListCommandHandler) Execute(ctx context.Context, req foodmodel.ListFoodReq) (*foodmodel.ListFoodRes, error) {
	foods, total, err := s.repo.ListFoods(ctx, req)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	foodDtos := convertListCategoryRes(foods)
	var resp foodmodel.ListFoodRes
	resp.Items = foodDtos
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}

func convertListCategoryRes(foods []foodmodel.Food) []foodmodel.FoodSearchResDto {
	var listCategoryRes []foodmodel.FoodSearchResDto
	for _, f := range foods {
		var listCatsDto foodmodel.FoodSearchResDto
		listCatsDto.Id = f.Id
		listCatsDto.Name = f.Name
		listCatsDto.Description = f.Description
		listCatsDto.Status = f.Status
		listCatsDto.UpdatedAt = f.UpdatedAt
		listCategoryRes = append(listCategoryRes, listCatsDto)
	}
	return listCategoryRes
}
