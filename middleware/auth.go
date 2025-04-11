package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func extractToken(authorizationStr string) (string, error) {
	token := strings.TrimPrefix(authorizationStr, "Bearer ")
	if token == "" {
		panic(errors.New("token is required"))
	}
	return token, nil
}

type ITokenValidator interface {
	Validate(token string) (datatype.Requester, error)
}

func Auth(tokenValidator ITokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c.GetHeader("Authorization"))
		if err != nil {
			panic(err)
		}

		fmt.Println("token: ", token)
		requester, err := tokenValidator.Validate(token)
		if err != nil {
			panic(err)
		}
		fmt.Println("requester: ", requester)

		c.Set(datatype.KeyRequester, requester)
		c.Next()
	}
}
