package grpcctrl

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/food"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
)

type FoodRepository interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]foodmodel.FoodInfoDto, error)
}

type UpdateService interface {
	Execute(ctx context.Context, req service.FoodUpdateReq) error
}

type FoodGrpcServer struct {
	food.UnimplementedFoodServer
	repo FoodRepository
}

func NewFoodGrpcServer(repo FoodRepository) *FoodGrpcServer {
	return &FoodGrpcServer{repo: repo}
}

func (f *FoodGrpcServer) GetFooodsByIds(ctx context.Context, req *food.GetFoodIdsRequest) (*food.FoodIdsResp, error) {
	uuidIds := make([]uuid.UUID, len(req.Ids))
	for i, id := range req.Ids {
		uuidIds[i] = uuid.MustParse(id)
	}

	cats, err := f.repo.FindByIds(ctx, uuidIds)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*food.FoodDTO, len(cats))

	for i, cat := range cats {
		result[i] = &food.FoodDTO{
			Id:           cat.Id.String(),
			Name:         cat.Name,
			Description:  cat.Description,
			Images:       cat.Images,
			Price:        float32(cat.Price),
			Avgpoint:     float32(cat.AvgPoint),
			CommentQty:   int64(cat.CommentQty),
			CategoryId:   cat.CategoryId.String(),
			RestaurantId: cat.CategoryId.String(),
			Status:       cat.Status,
		}
	}

	return &food.FoodIdsResp{Data: result}, nil
}

func (f *FoodGrpcServer) UpdateFoodById(ctx context.Context, req *food.UpdateFoodRequest) (*food.UpdateFoodResp, error) {
	// TODO:

	return &food.UpdateFoodResp{Id: req.Id}, nil
}
