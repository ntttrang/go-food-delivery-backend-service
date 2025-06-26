package middleware

import "github.com/gin-gonic/gin"

type MiddlewareProvider struct {
	tokenInstropecter ITokenValidator
}

func NewMiddlewareProvider(tokenInstropecter ITokenValidator) *MiddlewareProvider {
	return &MiddlewareProvider{
		tokenInstropecter: tokenInstropecter,
	}
}

func (p *MiddlewareProvider) Auth() gin.HandlerFunc {
	return Auth(p.tokenInstropecter)
}
