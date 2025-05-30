package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// UpdateOrderStateRequest represents the request to update order state
type UpdateOrderStateRequest struct {
	NewState      string  `json:"newState" binding:"required"`
	ShipperID     *string `json:"shipperId,omitempty"`
	PaymentStatus *string `json:"paymentStatus,omitempty"`
}

// AssignShipperRequest represents the request to assign a shipper
type AssignShipperRequest struct {
	ShipperID string `json:"shipperId" binding:"required"`
}

// UpdatePaymentStatusRequest represents the request to update payment status
type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"paymentStatus" binding:"required"`
}

// UpdateOrderStateAPI handles order state transitions
func (ctrl *OrderHttpController) UpdateOrderStateAPI(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		panic(datatype.ErrBadRequest.WithError("order ID is required"))
	}

	var req UpdateOrderStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Get user ID from requester context for audit
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	_ = requester.Subject().String() // updatedBy for future use

	// Call business logic in service
	// Note: This requires the order state management service to be available
	// For now, we'll use the update command handler
	updateReq := service.OrderUpdateReq{
		ID:            orderID,
		ShipperID:     req.ShipperID,
		State:         &req.NewState,
		PaymentStatus: req.PaymentStatus,
	}

	if err := ctrl.updateCmdHdl.Execute(c.Request.Context(), updateReq); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"orderId": orderID,
			"message": "Order state updated successfully",
		},
	})
}

// AssignShipperAPI handles shipper assignment to an order
func (ctrl *OrderHttpController) AssignShipperAPI(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		panic(datatype.ErrBadRequest.WithError("order ID is required"))
	}

	var req AssignShipperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Create update request with shipper assignment
	updateReq := service.OrderUpdateReq{
		ID:        orderID,
		ShipperID: &req.ShipperID,
	}

	if err := ctrl.updateCmdHdl.Execute(c.Request.Context(), updateReq); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"orderId":   orderID,
			"shipperId": req.ShipperID,
			"message":   "Shipper assigned successfully",
		},
	})
}

// UpdatePaymentStatusAPI handles payment status updates
func (ctrl *OrderHttpController) UpdatePaymentStatusAPI(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		panic(datatype.ErrBadRequest.WithError("order ID is required"))
	}

	var req UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Validate payment status
	if req.PaymentStatus != "pending" && req.PaymentStatus != "paid" {
		panic(datatype.ErrBadRequest.WithError("payment status must be 'pending' or 'paid'"))
	}

	// Create update request with payment status
	updateReq := service.OrderUpdateReq{
		ID:            orderID,
		PaymentStatus: &req.PaymentStatus,
	}

	if err := ctrl.updateCmdHdl.Execute(c.Request.Context(), updateReq); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"orderId":       orderID,
			"paymentStatus": req.PaymentStatus,
			"message":       "Payment status updated successfully",
		},
	})
}
