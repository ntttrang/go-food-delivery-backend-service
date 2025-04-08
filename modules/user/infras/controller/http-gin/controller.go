package httpgin

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
)

type IRegisterUserCommandHandler interface {
	Execute(ctx context.Context, req usermodel.RegisterUserReq) error
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

func (ctrl *UserHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("/register", ctrl.RegisterAPI)
	g.POST("/authenticate", ctrl.AuthenticateAPI)

	jwtComp := shareComponent.NewJwtComp(os.Getenv("JWT_SECRET_KEY"), 3600*24*7)
	g.GET("/profile", middleware.Auth(jwtComp), ctrl.GetProfileAPI)
}
