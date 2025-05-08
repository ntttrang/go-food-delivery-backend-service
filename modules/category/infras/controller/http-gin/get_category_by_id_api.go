package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CategoryHttpController) GetCategoryByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	category, err := ctrl.getDetailCmdHdl.Execute(c.Request.Context(), service.CategoryDetailReq{Id: id})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}
