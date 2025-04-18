package service

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IBulkCreateRepo interface {
	BulkInsert(ctx context.Context, data []categorymodel.Category) error
}

type BulkCreateCommandHandler struct {
	catRepo IBulkCreateRepo
}

func NewBulkCreateCommandHandler(catRepo IBulkCreateRepo) *BulkCreateCommandHandler {
	return &BulkCreateCommandHandler{catRepo: catRepo}
}

func (s *BulkCreateCommandHandler) Execute(ctx context.Context, datas []categorymodel.CategoryInsertDto) ([]uuid.UUID, error) {

	var categories []categorymodel.Category
	var ids []uuid.UUID
	for _, data := range datas {
		if err := data.Validate(); err != nil {
			return nil, datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
		}

		category := data.ConvertToCategory()
		category.Id, _ = uuid.NewV7()
		category.Status = sharedModel.StatusActive // Always set Active Status when insert

		categories = append(categories, *category)
		ids = append(ids, category.Id)
	}

	if err := s.catRepo.BulkInsert(ctx, categories); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	return ids, nil
}
