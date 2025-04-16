package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) CreateRestaurantCommentAPI(c *gin.Context) {
	var req restaurantmodel.RestaurantCommentCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserID = requester.Subject()

	if err := ctrl.createCommentRestaurantCmdHandler.Execute(c.Request.Context(), &req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": req.Id})
}
