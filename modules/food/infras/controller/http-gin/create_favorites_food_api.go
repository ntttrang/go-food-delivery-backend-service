package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) UpdateFavoritesFoodAPI(c *gin.Context) {
	var req foodmodel.FoodLike

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserId = requester.Subject()

	check, err := ctrl.addFavoritesCmdHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": check})
}
