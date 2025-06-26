package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"resty.dev/v3"
)

type CategoryDto struct {
	Id   uuid.UUID
	Name string
}

type CategoryRPCClient struct {
	catServiceURL string
}

func NewCategoryRPCClient(catServiceURL string) *CategoryRPCClient {
	return &CategoryRPCClient{catServiceURL: catServiceURL}
}

func (c *CategoryRPCClient) FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]CategoryDto, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data []CategoryDto `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/find-by-ids", c.catServiceURL)

	_, err := client.R().
		SetBody(map[string]interface{}{
			"ids": ids,
		}).
		SetResult(&response).
		Post(url)

	if err != nil {
		return nil, err
	}

	cats := response.Data
	catsMap := make(map[uuid.UUID]CategoryDto, len(cats))
	for _, r := range cats {
		catsMap[r.Id] = r
	}
	return catsMap, nil
}
