package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (ctrl *CategoryHttpController) CreateBulkCategoryAPI(c *gin.Context) {
	var requestBodyData []categorymodel.CategoryInsertDto

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call business logic in service
	ids, err := ctrl.bulkCreateCmdHdl.Execute(c.Request.Context(), requestBodyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"ids": ids}})
}
