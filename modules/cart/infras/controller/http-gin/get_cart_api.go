package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CartHttpController) GetCartByUserIdAndFoodIdAPI(c *gin.Context) {
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	foodId, err := uuid.Parse(c.Param("foodId"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	result, err := ctrl.getDetailQueryHdl.Execute(c.Request.Context(), service.CartDetailReq{
		UserID: userId,
		FoodID: foodId,
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
