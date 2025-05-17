package sharerpc

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
	"resty.dev/v3"
)

type IntrospectRpcClient struct {
	userServiceURL string
}

func NewIntrospectRpcClient(userServiceURL string) *IntrospectRpcClient {
	return &IntrospectRpcClient{
		userServiceURL: userServiceURL,
	}
}

type dataRequester struct {
	UserID    uuid.UUID `json:"id"`
	RoleValue string    `json:"role"`
}

func (r *dataRequester) Subject() uuid.UUID {
	return r.UserID
}

func (r *dataRequester) GetRole() string {
	return r.RoleValue
}

func (c *IntrospectRpcClient) Validate(token string) (datatype.Requester, error) {
	client := resty.New()

	type ResponseDTO struct {
		Data struct {
			UserId string `json:"id"`
			Role   string `json:"role"`
		} `json:"data"`
	}

	var response ResponseDTO

	url := fmt.Sprintf("%s/introspect-token", c.userServiceURL)

	_, err := client.R().SetBody(map[string]interface{}{
		"token": token,
	}).SetResult(&response).Post(url)
	if err != nil {
		fmt.Println(err)
		return nil, errors.WithStack(err)
	}
	fmt.Println(response)

	return &dataRequester{
		UserID:    uuid.MustParse(response.Data.UserId),
		RoleValue: response.Data.Role,
	}, nil
}
