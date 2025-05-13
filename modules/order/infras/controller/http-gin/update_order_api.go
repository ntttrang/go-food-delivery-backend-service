package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *OrderHttpController) UpdateOrderAPI(c *gin.Context) {
	id := c.Param("id")
	
	var req service.OrderUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Set ID from path parameter
	req.ID = id

	// Call business logic in service
	if err := ctrl.updateCmdHdl.Execute(c.Request.Context(), req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
