package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (ctrl *UserHttpController) RefreshAPI(c *gin.Context) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, datatype.ErrBadRequest.WithError(err.Error()))
		return
	}

	// Refresh the access token
	authRes, err := ctrl.refreshCmdHdl.Execute(c.Request.Context(), req.RefreshToken)
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
