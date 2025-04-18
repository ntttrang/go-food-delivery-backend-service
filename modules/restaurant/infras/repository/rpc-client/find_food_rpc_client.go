package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"resty.dev/v3"
)

func (c *FoodRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.Foods, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data []restaurantmodel.Foods `json:"data"`
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
	foodMap := make(map[uuid.UUID]restaurantmodel.Foods, len(foods))
	for _, r := range foods {
		foodMap[r.Id] = r
	}
	return foodMap, nil
}
