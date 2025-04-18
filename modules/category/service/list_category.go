package service

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IListRep interface {
	ListCategories(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.Category, int64, error)
}

type ListCommandHandler struct {
	catRepo IListRep
}

func NewListCommandHandler(catRepo IListRep) *ListCommandHandler {
	return &ListCommandHandler{
		catRepo: catRepo,
	}
}

func (s *ListCommandHandler) Execute(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.ListCategoryRes, int64, error) {
	categories, total, err := s.catRepo.ListCategories(ctx, req)

	if err != nil {
		return nil, 0, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return convertListCategoryRes(categories), total, nil
}

func convertListCategoryRes(cats []categorymodel.Category) []categorymodel.ListCategoryRes {
	var listCategoryRes []categorymodel.ListCategoryRes
	for _, cat := range cats {
		var listCatsDto categorymodel.ListCategoryRes
		listCatsDto.Id = cat.Id
		listCatsDto.Name = cat.Name
		listCatsDto.Description = cat.Description
		listCatsDto.Status = cat.Status
		listCatsDto.UpdatedAt = cat.UpdatedAt
		listCategoryRes = append(listCategoryRes, listCatsDto)
	}
	return listCategoryRes
}
