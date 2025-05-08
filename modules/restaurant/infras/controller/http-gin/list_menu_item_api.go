package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) ListMenuItemAPI(c *gin.Context) {
	restaurantId, err := uuid.Parse(c.Param("restaurantId"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	result, err := ctrl.listMenuItemQueryHandler.Execute(c.Request.Context(), restaurantId)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
