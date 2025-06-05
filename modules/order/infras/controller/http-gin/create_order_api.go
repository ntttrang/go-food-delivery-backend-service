package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// CreateOrderAPI is restricted to ADMIN users only
func (ctrl *OrderHttpController) CreateOrderAPI(c *gin.Context) {
	var req service.OrderCreateDto

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Get user ID and role from requester context
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)

	// Check if user has ADMIN role - only ADMIN can use this API
	if requester.GetRole() != string(datatype.RoleAdmin) {
		panic(datatype.ErrForbidden.WithError("only administrators can create orders manually"))
	}

	req.UserID = requester.Subject().String()

	// Call business logic in service
	orderId, err := ctrl.createCmdHdl.Execute(c.Request.Context(), &req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"orderId": orderId,
			"message": "Order created successfully by administrator",
		},
	})
}
