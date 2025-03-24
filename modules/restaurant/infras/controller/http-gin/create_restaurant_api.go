package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/internal/model"
)

func (ctl *RestaurantHttpController) CreateRestaurantAPI(c *gin.Context) {
	var requestBodyData restaurantmodel.RestaurantInsertDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call business logic in service
	if err := ctl.restaurantService.CreateRestaurant(c.Request.Context(), &requestBodyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
