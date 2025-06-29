package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"resty.dev/v3"
)

type CartRPCClient struct {
	cartServiceURL string
}

type CartSummaryData struct {
	CartID       uuid.UUID `json:"cartId"`
	UserID       uuid.UUID `json:"userId"`
	RestaurantID uuid.UUID `json:"restaurantId"`
	FoodId       uuid.UUID `json:"foodId"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	ItemCount    int       `json:"itemCount"`
	TotalPrice   float64   `json:"totalPrice"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

func NewCartRPCClient(cartServiceURL string) *CartRPCClient {
	return &CartRPCClient{cartServiceURL: cartServiceURL}
}

func (c *CartRPCClient) FindById(ctx context.Context, cartId string, userId string) ([]CartSummaryData, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data []CartSummaryData `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/cart-summary", c.cartServiceURL)

	_, err := client.R().
		SetQueryParam("cartId", cartId).
		SetQueryParam("userId", userId).
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
