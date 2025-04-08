package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (ctrl *RestaurantHttpController) CreateRestaurantAPI(c *gin.Context) {
	var requestBodyData restaurantmodel.RestaurantInsertDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), &requestBodyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
