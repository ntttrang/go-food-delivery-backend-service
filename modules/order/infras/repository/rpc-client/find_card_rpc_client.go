package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"resty.dev/v3"
)

type CardRPCClient struct {
	paymentServiceURL string
}

func NewCardRPCClient(paymentServiceURL string) *CardRPCClient {
	return &CardRPCClient{paymentServiceURL: paymentServiceURL}
}

func (c *CardRPCClient) FindById(ctx context.Context, id uuid.UUID) (*ordermodel.Card, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data ordermodel.Card `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/find-by-id", c.paymentServiceURL)

	_, err := client.R().
		SetBody(map[string]interface{}{
			"id": id,
		}).
		SetResult(&response).
		Post(url)

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
