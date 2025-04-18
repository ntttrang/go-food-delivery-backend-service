package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) GetRestaurantDetailAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	restaurant, err := ctrl.getDetailQueryHdl.Execute(c.Request.Context(), restaurantmodel.RestaurantDetailReq{Id: id})

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": restaurant})
}
