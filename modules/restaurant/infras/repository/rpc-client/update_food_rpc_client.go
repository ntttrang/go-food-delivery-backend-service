package rpcclient

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

type FoodUpdateReq struct {
	RestaurantId string `json:"restaurantId"` // Can be empty or missing if data type = string. Otherwise, uuid.UUID isn't
	CategoryId   string `json:"categoryId"`   // Can be empty or missing if data type = string. Otherwise, uuid.UUID isn't
	FoodId       string
}

type FoodUpdateRes struct {
	Id string `json:"id"`
}

type FoodRPCClient struct {
	foodServiceURL string
}

func NewFoodRPCClient(foodServiceURL string) *FoodRPCClient {
	return &FoodRPCClient{foodServiceURL: foodServiceURL}
}

func (c *FoodRPCClient) UpdateFoods(ctx context.Context, req FoodUpdateReq) (*string, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data FoodUpdateRes `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/update/%s", c.foodServiceURL, req.FoodId)
	fmt.Println("url:", url)
	_, err := client.R().
		SetBody(req).
		SetResult(&response).
		Patch(url)

	if err != nil {
		return nil, err
	}

	return &response.Data.Id, nil
}
