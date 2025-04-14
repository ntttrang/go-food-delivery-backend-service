package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) CreateUserAPI(c *gin.Context) {
	var requestBodyData *usermodel.CreateUserReq

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), requestBodyData); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
