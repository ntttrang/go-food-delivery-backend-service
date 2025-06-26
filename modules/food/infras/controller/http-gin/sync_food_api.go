package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// SyncFoodByIdAPI handles the sync food by ID API endpoint
func (ctrl *FoodHttpController) SyncFoodByIdAPI(c *gin.Context) {
	// Check if search functionality is available
	if ctrl.syncFoodIndexCommandHandler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Sync functionality is not available. Elasticsearch is not configured.",
		})
		return
	}

	// Get food ID from path
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError("Invalid food ID"))
	}

	// Sync food
	err = ctrl.syncFoodByIdCommandHandler.SyncFood(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": "Food indexed successfully"})
}
