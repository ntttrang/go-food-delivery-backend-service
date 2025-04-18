package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) DeleteRestaurantCommentAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := ctrl.deleteRestaurantCmdHdl.Execute(c.Request.Context(), restaurantmodel.RestaurantDeleteCommentReq{Id: id}); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}
