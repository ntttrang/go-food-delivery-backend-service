package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// SyncRestaurantByIdAPI handles the sync restaurant by ID API endpoint
func (ctrl *RestaurantHttpController) SyncRestaurantByIdAPI(c *gin.Context) {
	// Check if search functionality is available
	if ctrl.syncRestaurantIndexCommandHandler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Sync functionality is not available. Elasticsearch is not configured.",
		})
		return
	}

	// Get restaurant ID from path
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError("Invalid restaurant ID"))
	}

	// Sync restaurant
	err = ctrl.syncRestaurantIndexCommandHandler.SyncRestaurant(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": "Restaurant indexed successfully"})
}
