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
type IntrospectCommandHandler interface {
	Execute(ctx context.Context, req usermodel.IntrospectReq) (*usermodel.IntrospectRes, error)
}

type UserHttpController struct {
	registerUserCmdHdl IRegisterUserCommandHandler
	authCmdHdl         IAuthenticateCommandHandler
	introspectCmdHdl   IntrospectCommandHandler
}

func NewUserHttpController(registerUserCmdHdl IRegisterUserCommandHandler, authCmdHdl IAuthenticateCommandHandler, introspectCmdHdl IntrospectCommandHandler) *UserHttpController {
	return &UserHttpController{
		registerUserCmdHdl: registerUserCmdHdl,
		authCmdHdl:         authCmdHdl,
		introspectCmdHdl:   introspectCmdHdl,
	}
}

func (ctrl *UserHttpController) SetupRoutes(g *gin.RouterGroup, authMld gin.HandlerFunc) {
	g.POST("/register", ctrl.RegisterAPI)
	g.POST("/authenticate", ctrl.AuthenticateAPI) // Login
	g.GET("/profile", authMld, ctrl.GetProfileAPI)

	g.POST("/rpc/users/introspect-token", ctrl.IntrospectTokenRpcAPI)
}
