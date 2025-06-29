package grpcclient

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/food"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FoodGRPCClient struct {
	foodGRPCServerURL string
	client            food.FoodClient
}

func NewFoodGRPCClient(categoryGRPCServerURL string) *FoodGRPCClient {
	conn, err := grpc.NewClient(
		categoryGRPCServerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := food.NewFoodClient(conn)
	return &FoodGRPCClient{foodGRPCServerURL: categoryGRPCServerURL, client: client}
}

func (c *FoodGRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]ordermodel.Food, error) {
	strIds := make([]string, len(ids))

	for i, id := range ids {
		strIds[i] = id.String()
	}

	resp, err := c.client.GetFooodsByIds(ctx, &food.GetFoodIdsRequest{Ids: strIds})

	if err != nil {
		return nil, err
	}

	foods := resp.Data
	foodsMap := make(map[uuid.UUID]ordermodel.Food, len(foods))
	for _, r := range foods {
		uuidId := uuid.MustParse(r.Id)
		v := ordermodel.Food{
			Id:           uuidId,
			Name:         r.Name,
			Description:  r.Name,
			Images:       r.Images,
			Price:        float64(r.Price),
			AvgPoint:     float64(r.Avgpoint),
			CommentQty:   int(r.CommentQty),
			CategoryId:   uuid.MustParse(r.CategoryId),
			RestaurantId: uuid.MustParse(r.RestaurantId),
			Status:       r.Status,
		}
		foodsMap[uuidId] = v
	}
	return foodsMap, nil
}
