package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// CreateOrderFromCartAPI handles creating an order from a cart (for standard customer order creation)
func (ctrl *OrderHttpController) CreateOrderFromCartAPI(c *gin.Context) {
	var req service.OrderCreateFromCartDto

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Get user ID from requester context
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserID = requester.Subject().String()

	// Call business logic in service
	orderId, err := ctrl.createFromCartCmdHdl.ExecuteFromCart(c.Request.Context(), &req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"orderId": orderId,
			"message": "Order created successfully from cart",
		},
	})
}
