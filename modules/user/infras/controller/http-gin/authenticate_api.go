package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) AuthenticateAPI(c *gin.Context) {
	var req service.AuthenticateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, datatype.ErrBadRequest.WithError(err.Error()))
		return
	}

	authRes, err := ctrl.authCmdHdl.Execute(c.Request.Context(), req)
	if err != nil {
		// Handle application errors with proper status codes
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			c.JSON(appErr.StatusCode(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, datatype.ErrInternalServerError.WithDebug(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, datatype.ResponseSuccess(authRes))
}
