package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) RPCUpdateFoodByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	var req foodmodel.FoodUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	req.Id = id

	if err := ctrl.updateByIdCommandHandler.Execute(c.Request.Context(), req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id}})
}
