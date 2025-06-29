package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
)

func (ctrl *OrderHttpController) DeleteOrderAPI(c *gin.Context) {
	id := c.Param("id")
	
	req := service.OrderDeleteReq{
		ID: id,
	}

	// Call business logic in service
	if err := ctrl.deleteCmdHdl.Execute(c.Request.Context(), req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
