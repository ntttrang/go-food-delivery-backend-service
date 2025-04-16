package service

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
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

func (s *ListCommandHandler) Execute(ctx context.Context, req foodmodel.ListFoodReq) ([]foodmodel.ListFoodRes, int64, error) {
	categories, total, err := s.repo.ListFoods(ctx, req)

	if err != nil {
		return nil, 0, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return convertListCategoryRes(categories), total, nil
}

func convertListCategoryRes(cats []foodmodel.Food) []foodmodel.ListFoodRes {
	var listCategoryRes []foodmodel.ListFoodRes
	for _, cat := range cats {
		var listCatsDto foodmodel.ListFoodRes
		listCatsDto.Id = cat.Id
		listCatsDto.Name = cat.Name
		listCatsDto.Description = cat.Description
		listCatsDto.Status = cat.Status
		listCatsDto.UpdatedAt = cat.UpdatedAt
		listCategoryRes = append(listCategoryRes, listCatsDto)
	}
	return listCategoryRes
}
