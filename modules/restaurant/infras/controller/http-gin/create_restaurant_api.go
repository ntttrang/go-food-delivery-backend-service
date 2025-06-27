package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) CreateRestaurantAPI(c *gin.Context) {
	var requestBodyData *restaurantservice.RestaurantInsertDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	requestBodyData.OwnerId = requester.Subject()

	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), requestBodyData); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
