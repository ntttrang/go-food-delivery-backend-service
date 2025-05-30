package sharerpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"resty.dev/v3"
)

// CartRPCClient handles RPC communication with cart service
type CartRPCClient struct {
	cartServiceURL string
}

// NewCartRPCClient creates a new cart RPC client
func NewCartRPCClient(cartServiceURL string) *CartRPCClient {
	return &CartRPCClient{
		cartServiceURL: cartServiceURL,
	}
}

// CartItemRPCDto represents a cart item from RPC response
type CartItemRPCDto struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	FoodID       uuid.UUID `json:"foodId"`
	RestaurantID uuid.UUID `json:"restaurantId"`
	Quantity     int       `json:"quantity"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// CartItemsRPCResponse represents the RPC response structure
type CartItemsRPCResponse struct {
	Data []CartItemRPCDto `json:"data"`
}

// FindCartItemsByCartID retrieves cart items by cart ID via RPC
func (c *CartRPCClient) FindCartItemsByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) ([]CartItemRPCDto, error) {
	if c.cartServiceURL == "" {
		return nil, errors.New("cart service URL not configured")
	}

	client := resty.New()

	var response CartItemsRPCResponse

	// Construct the URL for getting cart items
	url := fmt.Sprintf("%s/carts/items", c.cartServiceURL)

	resp, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"cartId": cartID.String(),
			"userId": userID.String(),
		}).
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Check for HTTP errors
	if resp.StatusCode() >= 400 {
		return nil, errors.Errorf("cart service returned error: %d - %s", resp.StatusCode(), resp.String())
	}

	return response.Data, nil
}

// UpdateCartStatus updates the cart status via RPC
func (c *CartRPCClient) UpdateCartStatus(ctx context.Context, cartID uuid.UUID, status string) error {
	if c.cartServiceURL == "" {
		return errors.New("cart service URL not configured")
	}

	client := resty.New()

	// Construct the URL for updating cart status
	url := fmt.Sprintf("%s/carts/%s/status", c.cartServiceURL, cartID.String())

	requestBody := map[string]any{
		"status": status,
	}

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Patch(url)

	if err != nil {
		return errors.WithStack(err)
	}

	// Check for HTTP errors
	if resp.StatusCode() >= 400 {
		return errors.Errorf("cart service returned error: %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// ValidateCartOwnership validates that the cart belongs to the user
func (c *CartRPCClient) ValidateCartOwnership(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) error {
	if c.cartServiceURL == "" {
		return errors.New("cart service URL not configured")
	}

	client := resty.New()

	type CartValidationResponse struct {
		Data struct {
			Valid  bool   `json:"valid"`
			Reason string `json:"reason,omitempty"`
		} `json:"data"`
	}

	var response CartValidationResponse

	// Construct the URL for validating cart ownership
	url := fmt.Sprintf("%s/carts/%s/validate", c.cartServiceURL, cartID.String())

	resp, err := client.R().
		SetContext(ctx).
		SetQueryParam("userId", userID.String()).
		SetResult(&response).
		Get(url)

	if err != nil {
		return errors.WithStack(err)
	}

	// Check for HTTP errors
	if resp.StatusCode() >= 400 {
		return errors.Errorf("cart service returned error: %d - %s", resp.StatusCode(), resp.String())
	}

	if !response.Data.Valid {
		reason := response.Data.Reason
		if reason == "" {
			reason = "cart validation failed"
		}
		return errors.New(reason)
	}

	return nil
}

// CartSummaryRPCDto represents cart summary information
type CartSummaryRPCDto struct {
	CartID       uuid.UUID `json:"cartId"`
	UserID       uuid.UUID `json:"userId"`
	RestaurantID uuid.UUID `json:"restaurantId"`
	ItemCount    int       `json:"itemCount"`
	TotalPrice   float64   `json:"totalPrice"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// GetCartSummary gets cart summary information
func (c *CartRPCClient) GetCartSummary(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*CartSummaryRPCDto, error) {
	if c.cartServiceURL == "" {
		return nil, errors.New("cart service URL not configured")
	}

	client := resty.New()

	type CartSummaryResponse struct {
		Data CartSummaryRPCDto `json:"data"`
	}

	var response CartSummaryResponse

	// Construct the URL for getting cart summary
	url := fmt.Sprintf("%s/carts/%s/summary", c.cartServiceURL, cartID.String())

	resp, err := client.R().
		SetContext(ctx).
		SetQueryParam("userId", userID.String()).
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Check for HTTP errors
	if resp.StatusCode() >= 400 {
		return nil, errors.Errorf("cart service returned error: %d - %s", resp.StatusCode(), resp.String())
	}

	return &response.Data, nil
}
