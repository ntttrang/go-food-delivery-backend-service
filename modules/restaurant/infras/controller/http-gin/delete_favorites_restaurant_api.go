package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

type DeleteRequest struct {
	RestaurantID string `form:"restaurantId" binding:"required,uuid"`
	UserID       string `form:"userId" binding:"required,uuid"`
}

func (ctrl *RestaurantHttpController) DeleteFavoritesRestaurantAPI(c *gin.Context) {
	var delReq DeleteRequest
	if err := c.ShouldBindQuery(&delReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req restaurantmodel.RestaurantLike
	copier.Copy(&req, &delReq)

	if err := ctrl.deleteFavoritesCmdHdl.Execute(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": req})
}
