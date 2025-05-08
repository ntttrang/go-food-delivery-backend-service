package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) CreateMenuItemAPI(c *gin.Context) {
	var req restaurantservice.MenuItemCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := ctrl.createMenuItemCmdHdl.Execute(c.Request.Context(), &req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": req})
}
