package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CartHttpController) UpdateCartStatusAPI(c *gin.Context) {
	cartId, err := uuid.Parse(c.Query("cartId"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	status := c.Query("status")
	// Validate status values
	validStatuses := map[string]bool{
		"ACTIVE":    true,
		"UPDATED":   true,
		"PROCESSED": true,
	}

	if !validStatuses[status] {
		panic(datatype.ErrBadRequest.WithError("invalid status value"))
	}

	if err := ctrl.repo.UpdateCartStatusByCartID(c.Request.Context(), cartId, status); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": cartId}})
}
