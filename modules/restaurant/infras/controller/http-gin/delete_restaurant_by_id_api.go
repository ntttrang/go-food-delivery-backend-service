package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) DeleteRestaurantByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := ctrl.deleteCmdHdl.Execute(c.Request.Context(), restaurantservice.RestaurantDeleteReq{Id: id}); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": id})
}
