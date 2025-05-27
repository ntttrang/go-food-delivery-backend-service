package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"resty.dev/v3"
)

type RestaurantRPCClient struct {
	restaurantServiceURL string
}

type RPCGetByIdsResponseDTO struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Addr             string    `json:"addr"`
	CityId           int       `json:"cityId"`
	Lat              float64   `json:"lat"`
	Lng              float64   `json:"lng"`
	Cover            string    `json:"cover"`
	Logo             string    `json:"logo"`
	ShippingFeePerKm float64   `json:"shippingFeePerKm"`
}

func NewRestaurantRPCClient(restaurantServiceURL string) *RestaurantRPCClient {
	return &RestaurantRPCClient{restaurantServiceURL: restaurantServiceURL}
}

func (c *RestaurantRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]RPCGetByIdsResponseDTO, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data []RPCGetByIdsResponseDTO `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/find-by-ids", c.restaurantServiceURL)

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
	foodMap := make(map[uuid.UUID]RPCGetByIdsResponseDTO, len(foods))
	for _, r := range foods {
		foodMap[r.Id] = r
	}
	return foodMap, nil
}
