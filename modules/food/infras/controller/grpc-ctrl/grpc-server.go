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
	repo          FoodRepository
	updateService UpdateService
}

func NewFoodGrpcServer(repo FoodRepository, updateService UpdateService) *FoodGrpcServer {
	return &FoodGrpcServer{
		repo:          repo,
		updateService: updateService,
	}
}

func (f *FoodGrpcServer) GetFooodsByIds(ctx context.Context, req *food.GetFoodIdsRequest) (*food.FoodIdsResp, error) {
	log.Println("[START] GRPC - GetFooodsByIds")
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
			RestaurantId: cat.RestaurantId.String(),
			Status:       cat.Status,
		}
	}
	log.Println("[END] GRPC - GetFooodsByIds")
	return &food.FoodIdsResp{Data: result}, nil
}

func (f *FoodGrpcServer) UpdateFoodById(ctx context.Context, req *food.UpdateFoodRequest) (*food.UpdateFoodResp, error) {
	log.Println("[START] GRPC - UpdateFoodById")

	// Parse the food ID
	foodId, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Invalid food ID: %v", err)
		return nil, err
	}

	// Convert gRPC request to service request
	updateReq := service.FoodUpdateReq{
		Id:          foodId,
		Name:        req.Name,
		Description: req.Description,
	}

	// Handle optional fields (use pointers for optional fields)
	if req.Status != "" {
		updateReq.Status = &req.Status
	}
	if req.RestaurantId != "" {
		updateReq.RestaurantId = &req.RestaurantId
	}
	if req.CategoryId != "" {
		updateReq.CategoryId = &req.CategoryId
	}
	if req.Image != "" {
		updateReq.Image = &req.Image
	}

	// Execute the update service
	if err := f.updateService.Execute(ctx, updateReq); err != nil {
		log.Printf("Failed to update food: %v", err)
		return nil, err
	}
	log.Println("[END] GRPC - UpdateFoodById")
	return &food.UpdateFoodResp{Id: req.Id}, nil
}
