package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) GetProfileAPI(c *gin.Context) {
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)

	c.JSON(http.StatusOK, gin.H{"data": requester.Subject()})
}
