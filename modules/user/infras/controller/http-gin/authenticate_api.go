package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/utils"
)

func (ctrl *UserHttpController) AuthenticateAPI(c *gin.Context) {
	var req service.AuthenticateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, datatype.ErrBadRequest.WithError(err.Error()))
		return
	}

	// Auto-detect OS if not provided by client
	if req.OS == "" {
		userAgent := c.GetHeader("User-Agent")
		req.OS = utils.GetSimpleOS(userAgent)
	}

	// Set production flag if not provided
	if req.IsProduction == nil {
		production := true // default to production
		req.IsProduction = &production
	}

	// Get User-Agent for OS detection
	userAgent := c.GetHeader("User-Agent")

	authRes, err := ctrl.authCmdHdl.Execute(c.Request.Context(), req, userAgent)
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
