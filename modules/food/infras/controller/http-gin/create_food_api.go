package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) CreateFoodAPI(c *gin.Context) {
	var req service.FoodInsertDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, datatype.ErrBadRequest.WithError(err.Error()))
		return
	}

	// call business logic in service
	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), &req); err != nil {
		// Handle application errors with proper status codes
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			c.JSON(appErr.StatusCode(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, datatype.ErrInternalServerError.WithDebug(err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": req.Id})
}
