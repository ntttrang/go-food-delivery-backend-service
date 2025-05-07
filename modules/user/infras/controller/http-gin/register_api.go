package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) RegisterAPI(c *gin.Context) {
	var req service.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := ctrl.registerUserCmdHdl.Execute(c.Request.Context(), &req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, datatype.ResponseSuccess(gin.H{"id": req.Id}))
}
