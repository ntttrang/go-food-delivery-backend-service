package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
)

func (ctrl *OrderHttpController) ListOrdersAPI(c *gin.Context) {
	var req service.OrderListReq

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call business logic in service
	result, err := ctrl.listQueryHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, result)
}
