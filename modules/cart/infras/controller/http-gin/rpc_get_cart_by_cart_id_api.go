package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// GetCartSummaryAPI gets cart summary information
// This endpoint is used by the order service via RPC
// Calls repository directly for simplicity
func (ctrl *CartHttpController) GetCartSummaryAPI(c *gin.Context) {
	cartIDStr := c.Query("cartId")
	userIDStr := c.Query("userId")

	if userIDStr == "" {
		panic(datatype.ErrBadRequest.WithError("userId is required"))
	}

	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		panic(datatype.ErrBadRequest.WithError("invalid cartId format"))
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		panic(datatype.ErrBadRequest.WithError("invalid userId format"))
	}

	// Call repository directly to get cart summary
	summaries, err := ctrl.repo.GetCartSummaryByCartID(c.Request.Context(), cartID, userID)
	if err != nil {
		panic(datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summaries,
	})
}
