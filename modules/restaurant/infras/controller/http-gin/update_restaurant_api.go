package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) UpdateRestaurantByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	var requestBodyData restaurantservice.RestaurantUpdateDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	cmd := restaurantservice.RestaurantUpdateReq{
		Id:  id,
		Dto: requestBodyData,
	}

	if err := ctrl.updateCmdHdl.Execute(c.Request.Context(), cmd); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
