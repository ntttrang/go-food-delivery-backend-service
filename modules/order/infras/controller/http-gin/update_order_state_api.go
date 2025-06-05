package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// UpdateOrderStateRequest represents the unified request to update order state
// This single endpoint can handle:
// 1. State transitions (preparing, on_the_way, delivered, cancel)
// 2. Shipper assignment (via shipperId field)
// 3. Payment status updates (via paymentStatus field)
// 4. Order cancellation (via cancellationReason field)
type UpdateOrderStateRequest struct {
	NewState           string  `json:"newState" binding:"required"`  // Required: target state
	ShipperID          *string `json:"shipperId,omitempty"`          // Optional: assign shipper
	PaymentStatus      *string `json:"paymentStatus,omitempty"`      // Optional: update payment status
	CancellationReason *string `json:"cancellationReason,omitempty"` // Required when newState is "cancel"
}

// UpdateOrderStateAPI handles order state transitions, assign shipper, update payment status and cancel an order
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
	updatedBy := requester.Subject().String()

	// Create state transition request
	updateReq := &service.StateTransitionRequest{
		OrderID:            orderID,
		NewState:           req.NewState,
		ShipperID:          req.ShipperID,
		PaymentStatus:      req.PaymentStatus,
		UpdatedBy:          updatedBy,
		CancellationReason: req.CancellationReason,
	}

	if err := ctrl.updateOrderStateCmdHdl.Execute(c.Request.Context(), updateReq); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"orderId": orderID,
			"message": "Order state updated successfully",
		},
	})
}
