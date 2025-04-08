package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (ctrl *CategoryHttpController) CreateCategoryAPI(c *gin.Context) {
	var requestBodyData categorymodel.CategoryInsertDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call business logic in service
	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), &requestBodyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
