package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"resty.dev/v3"
)

func (c *CartRPCClient) UpdateCartStatus(ctx context.Context, cartID uuid.UUID, status string) error {
	client := resty.New()

	type ResponseDTO struct {
		Data ordermodel.Card `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/update-status", c.cartServiceURL)

	_, err := client.R().
		SetQueryParam("cartId", cartID.String()).
		SetQueryParam("status", status).
		SetResult(&response).
		Patch(url)

	if err != nil {
		return err
	}

	return nil
}
