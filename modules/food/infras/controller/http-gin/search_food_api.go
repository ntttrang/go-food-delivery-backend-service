package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) SearchFoodAPI(c *gin.Context) {
	var req foodmodel.FoodSearchReq
	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Process pagination
	req.Process()

	// Execute search
	result, err := ctrl.searchFoodQueryHandler.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (ctrl *FoodHttpController) ReindexFoodAPI(c *gin.Context) {
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
