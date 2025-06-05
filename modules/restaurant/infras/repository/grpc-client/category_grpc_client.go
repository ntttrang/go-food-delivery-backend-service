package grpcclient

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/category"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryGRPCClient struct {
	catGRPCServerURL string
	client           category.CategoryClient
}

func NewCategoryGRPCClient(categoryGRPCServerURL string) *CategoryGRPCClient {
	conn, err := grpc.NewClient(
		categoryGRPCServerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := category.NewCategoryClient(conn)
	return &CategoryGRPCClient{catGRPCServerURL: categoryGRPCServerURL, client: client}
}

func (c *CategoryGRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantservice.CategoryDto, error) {
	strIds := make([]string, len(ids))

	for i, id := range ids {
		strIds[i] = id.String()
	}

	resp, err := c.client.GetCategoriesByIds(ctx, &category.GetCatIdsRequest{Ids: strIds})

	if err != nil {
		return nil, err
	}

	cats := resp.Data
	catsMap := make(map[uuid.UUID]restaurantservice.CategoryDto, len(cats))
	for _, r := range cats {
		uuidId := uuid.MustParse(r.Id)
		v := restaurantservice.CategoryDto{
			Id:   uuidId,
			Name: r.Name,
		}
		catsMap[uuidId] = v
	}
	return catsMap, nil
}
