package service

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/model"
)

func (s *CategoryService) ListCategories(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.ListCategoryRes, int64, error) {
	categories, total, err := s.catRepo.ListCategories(ctx, req)

	if err != nil {
		return nil, 0, err
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
