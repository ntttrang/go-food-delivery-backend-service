package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) GetFoodByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	Food, err := ctrl.getDetailCmdHdl.Execute(c.Request.Context(), foodmodel.FoodDetailReq{Id: id})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": Food})
}
