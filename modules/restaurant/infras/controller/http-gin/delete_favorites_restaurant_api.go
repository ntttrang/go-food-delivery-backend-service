package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type DeleteRequest struct {
	RestaurantID string `form:"restaurantId" binding:"required,uuid"`
	UserID       string `form:"userId" binding:"required,uuid"`
}

func (ctrl *RestaurantHttpController) DeleteFavoritesRestaurantAPI(c *gin.Context) {
	var delReq DeleteRequest
	if err := c.ShouldBindQuery(&delReq); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	var req restaurantmodel.RestaurantLike
	copier.Copy(&req, &delReq)

	if err := ctrl.deleteFavoritesCmdHdl.Execute(c.Request.Context(), req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": req})
}
