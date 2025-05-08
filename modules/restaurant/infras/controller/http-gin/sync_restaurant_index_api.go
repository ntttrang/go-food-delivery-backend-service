package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SyncRestaurantIndexAPI handles the sync restaurant index API endpoint
func (ctrl *RestaurantHttpController) SyncRestaurantIndexAPI(c *gin.Context) {
	// Check if search functionality is available
	if ctrl.syncRestaurantIndexCommandHandler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Sync functionality is not available. Elasticsearch is not configured.",
		})
		return
	}

	// Sync all restaurants
	err := ctrl.syncRestaurantIndexCommandHandler.SyncAll(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": "Restaurant index synchronized successfully"})
}
