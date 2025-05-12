package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CartHttpController) ListCartItemAPI(c *gin.Context) {
	var searchDto service.CartItemSearchDto
	if err := c.ShouldBind(&searchDto); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	req := service.CartItemListReq{
		CartItemSearchDto: searchDto,
	}

	result, err := ctrl.listCartItemQueryHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result.Items,
	})
}
