package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CategoryHttpController) CreateCategoryAPI(c *gin.Context) {
	var requestBodyData service.CategoryInsertDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// call business logic in service
	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), &requestBodyData); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
