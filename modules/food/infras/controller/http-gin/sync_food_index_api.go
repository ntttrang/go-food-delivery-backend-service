package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SyncFoodIndexAPI handles the sync food index API endpoint
func (ctrl *FoodHttpController) SyncFoodIndexAPI(c *gin.Context) {
	// Check if search functionality is available
	if ctrl.syncFoodIndexCommandHandler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Sync functionality is not available. Elasticsearch is not configured.",
		})
		return
	}

	// This endpoint should be protected and only accessible by admins
	err := ctrl.syncFoodIndexCommandHandler.SyncAll(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food index synchronized successfully"})
}
