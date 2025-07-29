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
		c.JSON(http.StatusBadRequest, datatype.ErrBadRequest.WithError(err.Error()))
		return
	}

	// Get user ID and role from requester context
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)

	// Check if user has ADMIN role - only ADMIN can use this API
	if requester.GetRole() != string(datatype.RoleAdmin) {
		c.JSON(http.StatusForbidden, datatype.ErrForbidden.WithError("only administrators can create orders manually"))
		return
	}

	req.UserID = requester.Subject().String()

	// Call business logic in service
	orderId, err := ctrl.createCmdHdl.Execute(c.Request.Context(), &req)
	if err != nil {
		// Handle application errors with proper status codes
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			c.JSON(appErr.StatusCode(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, datatype.ErrInternalServerError.WithDebug(err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"orderId": orderId,
			"message": "Order created successfully by administrator",
		},
	})
}
