package service

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type ListCategoryReq struct {
	SearchCategoryDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

type ListCategoryRes struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Status      string    `json:"status"`
	sharedModel.DateDto
}

type SearchCategoryDto struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}

// Initilize service
type IListRep interface {
	ListCategories(ctx context.Context, req ListCategoryReq) ([]categorymodel.Category, int64, error)
}

type ListCommandHandler struct {
	catRepo IListRep
}

func NewListCommandHandler(catRepo IListRep) *ListCommandHandler {
	return &ListCommandHandler{
		catRepo: catRepo,
	}
}

// Implement
func (s *ListCommandHandler) Execute(ctx context.Context, req ListCategoryReq) ([]ListCategoryRes, int64, error) {
	categories, total, err := s.catRepo.ListCategories(ctx, req)

	if err != nil {
		return nil, 0, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return convertListCategoryRes(categories), total, nil
}

func convertListCategoryRes(cats []categorymodel.Category) []ListCategoryRes {
	var listCategoryRes []ListCategoryRes
	for _, cat := range cats {
		var listCatsDto ListCategoryRes
		listCatsDto.Id = cat.Id
		listCatsDto.Name = cat.Name
		listCatsDto.Description = cat.Description
		listCatsDto.Icon = cat.Icon
		listCatsDto.Status = cat.Status
		listCatsDto.CreatedAt = cat.CreatedAt
		listCatsDto.UpdatedAt = cat.UpdatedAt
		listCategoryRes = append(listCategoryRes, listCatsDto)
	}
	return listCategoryRes
}
