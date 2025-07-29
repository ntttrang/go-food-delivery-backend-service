package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type LogoutReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (ctrl *UserHttpController) LogoutAPI(c *gin.Context) {
	var req LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, datatype.ErrBadRequest.WithError(err.Error()))
		return
	}

	// Revoke the refresh token
	if err := ctrl.logoutCmdHdl.Execute(c.Request.Context(), req.RefreshToken); err != nil {
		// Handle application errors with proper status codes
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			c.JSON(appErr.StatusCode(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, datatype.ErrInternalServerError.WithDebug(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, datatype.ResponseSuccess(gin.H{"message": "Successfully logged out"}))
}
