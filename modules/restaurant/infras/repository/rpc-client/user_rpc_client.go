package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"resty.dev/v3"
)

type UserRPCClient struct {
	userServiceURL string
}

func NewUserRPCClient(userServiceURL string) *UserRPCClient {
	return &UserRPCClient{userServiceURL: userServiceURL}
}

func (c *UserRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.User, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data []restaurantmodel.User `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/find-by-ids", c.userServiceURL)

	_, err := client.R().
		SetBody(map[string]interface{}{
			"ids": ids,
		}).
		SetResult(&response).
		Post(url)

	if err != nil {
		return nil, err
	}

	users := response.Data
	userMap := make(map[uuid.UUID]restaurantmodel.User, len(users))
	for _, r := range users {
		userMap[r.Id] = r
	}
	return userMap, nil
}
