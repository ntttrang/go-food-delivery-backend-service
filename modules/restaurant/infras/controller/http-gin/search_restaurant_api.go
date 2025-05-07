package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// SearchRestaurantsAPI handles the search restaurants API endpoint
func (ctrl *RestaurantHttpController) SearchRestaurantsAPI(c *gin.Context) {
	// Create a request object to bind the JSON body
	var req restaurantmodel.RestaurantSearchReq

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Process pagination
	req.Process()

	// Execute search
	result, err := ctrl.searchRestaurantQueryHandler.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
