package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
)

func (ctrl *OrderHttpController) GetOrderDetailAPI(c *gin.Context) {
	id := c.Param("id")
	
	req := service.OrderDetailReq{
		ID: id,
	}

	// Call business logic in service
	result, err := ctrl.getDetailQueryHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, result)
}
