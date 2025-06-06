package grpcctrl

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/category"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

type CategoryRepository interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]categorymodel.Category, error)
}

type CategoryGrpcServer struct {
	category.UnimplementedCategoryServer
	repo CategoryRepository
}

func NewCategoryGrpcServer(repo CategoryRepository) *CategoryGrpcServer {
	return &CategoryGrpcServer{repo: repo}
}

func (s *CategoryGrpcServer) GetCategoriesByIds(ctx context.Context, req *category.GetCatIdsRequest) (*category.CatIdsResp, error) {
	log.Println("[START] GRPC - GetCategoriesByIds")

	uuidIds := make([]uuid.UUID, len(req.Ids))
	for i, id := range req.Ids {
		uuidIds[i] = uuid.MustParse(id)
	}

	cats, err := s.repo.FindByIds(ctx, uuidIds)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*category.CategoryDTO, len(cats))

	for i, cat := range cats {
		result[i] = &category.CategoryDTO{
			Id:          cat.Id.String(),
			Name:        cat.Name,
			Description: cat.Description,
			Icon:        cat.Icon,
			Status:      cat.Status,
		}
	}

	log.Println("[END] GRPC - GetCategoriesByIds")
	return &category.CatIdsResp{Data: result}, nil
}
