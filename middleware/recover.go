package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type CanGetStatusCode interface {
	StatusCode() int
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {

			isProduction := os.Getenv("ENV") == "prod" || os.Getenv("GIN_MODE") == "release"

			if r := recover(); r != nil {
				if appErr, ok := r.(CanGetStatusCode); ok {
					c.JSON(appErr.StatusCode(), appErr)
					if !isProduction {
						log.Printf("Error: %+v", appErr)
						panic(r)
					}
					return
				}

				appError := datatype.ErrInternalServerError

				if isProduction {
					c.JSON(appError.StatusCode(), appError.WithDebug(""))
				} else {
					c.JSON(appError.StatusCode(), appError.WithDebug(fmt.Sprintf("%s", r)))
					panic(r)
				}
			}
		}()

		c.Next()

	}
}
