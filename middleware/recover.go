package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong",
				})
				panic(r)
			}
		}()

		c.Next()
	}
}
