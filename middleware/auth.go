package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func extractToken(authorizationStr string) (string, error) {
	if authorizationStr == "" {
		return "", datatype.ErrUnauthorized.WithError("authorization header is required")
	}

	token := strings.TrimPrefix(authorizationStr, "Bearer ")
	if token == "" || token == authorizationStr {
		return "", datatype.ErrUnauthorized.WithError("bearer token is required")
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
			c.JSON(err.(interface{ StatusCode() int }).StatusCode(), err)
			c.Abort()
			return
		}

		requester, err := tokenValidator.Validate(token)
		if err != nil {
			// Handle token validation errors
			if appErr, ok := err.(interface{ StatusCode() int }); ok {
				c.JSON(appErr.StatusCode(), appErr)
			} else {
				c.JSON(datatype.ErrUnauthorized.StatusCode(), datatype.ErrUnauthorized.WithError(err.Error()))
			}
			c.Abort()
			return
		}

		c.Set(datatype.KeyRequester, requester)
		c.Next()
	}
}
