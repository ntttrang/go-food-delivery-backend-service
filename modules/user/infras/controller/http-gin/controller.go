package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
)

type IRegisterUserCommandHandler interface {
	Execute(ctx context.Context, req *usermodel.RegisterUserReq) error
}

type IAuthenticateCommandHandler interface {
	Execute(ctx context.Context, req usermodel.AuthenticateReq) (*usermodel.AuthenticateRes, error)
}

type UserHttpController struct {
	registerUserCmdHdl IRegisterUserCommandHandler
	authCmdHdl         IAuthenticateCommandHandler
}

func NewUserHttpController(registerUserCmdHdl IRegisterUserCommandHandler, authCmdHdl IAuthenticateCommandHandler) *UserHttpController {
	return &UserHttpController{
		registerUserCmdHdl: registerUserCmdHdl,
		authCmdHdl:         authCmdHdl,
	}
}

func (ctrl *UserHttpController) SetupRoutes(g *gin.RouterGroup, authMld gin.HandlerFunc) {
	g.POST("/register", ctrl.RegisterAPI)
	g.POST("/authenticate", ctrl.AuthenticateAPI)

	g.GET("/profile", authMld, ctrl.GetProfileAPI)
}
