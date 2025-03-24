package service

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/model"
)

type ICategoryRepository interface {
	Insert(ctx context.Context, data *categorymodel.Category) error
	BulkInsert(ctx context.Context, data []categorymodel.Category) error
	ListCategories(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.Category, int64, error)
}

type CategoryService struct {
	catRepo ICategoryRepository
}

func NewCategoryService(catRepo ICategoryRepository) *CategoryService {
	return &CategoryService{catRepo: catRepo}
}
