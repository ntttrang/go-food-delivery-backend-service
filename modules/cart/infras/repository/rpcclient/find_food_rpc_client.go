package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"resty.dev/v3"
)

type FoodRPCClient struct {
	foodServiceURL string
}

func NewFoodRPCClient(foodServiceURL string) *FoodRPCClient {
	return &FoodRPCClient{foodServiceURL: foodServiceURL}
}

func (c *FoodRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]cartmodel.Food, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data []cartmodel.Food `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/find-by-ids", c.foodServiceURL)

	_, err := client.R().
		SetBody(map[string]interface{}{
			"ids": ids,
		}).
		SetResult(&response).
		Post(url)

	if err != nil {
		return nil, err
	}

	foods := response.Data
	foodMap := make(map[uuid.UUID]cartmodel.Food, len(foods))
	for _, r := range foods {
		foodMap[r.Id] = r
	}
	return foodMap, nil
}
