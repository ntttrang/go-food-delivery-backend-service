package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CartHttpController) DeleteCartByIdAPI(c *gin.Context) {
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	foodId, err := uuid.Parse(c.Param("foodId"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := ctrl.deleteCmdHdl.Execute(c.Request.Context(), service.CartDeleteReq{UserID: userId, FoodID: foodId}); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}
