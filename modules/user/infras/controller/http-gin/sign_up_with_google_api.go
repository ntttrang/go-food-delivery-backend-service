package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *UserHttpController) SignUpWithGoogleAPI(c *gin.Context) {
	url := ctrl.signUpGgCmdHdl.GetAuthCodeUrl(c.Request.Context())
	c.Redirect(http.StatusTemporaryRedirect, url)
}
