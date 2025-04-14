package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type IGenerateCode interface {
	Execute(ctx context.Context, userId uuid.UUID) (string, error)
}

type IVerifyCode interface {
	Execute(ctx context.Context, userId uuid.UUID, code string) (bool, error)
}

type UserHttpController struct {
	registerUserCmdHdl IRegisterUserCommandHandler
	authCmdHdl         IAuthenticateCommandHandler
	introspectCmdHdl   IntrospectCommandHandler
	generateCode       IGenerateCode
	verifyCode         IVerifyCode
}

func NewUserHttpController(registerUserCmdHdl IRegisterUserCommandHandler, authCmdHdl IAuthenticateCommandHandler, introspectCmdHdl IntrospectCommandHandler,
	generateCode IGenerateCode, verifyCode IVerifyCode) *UserHttpController {
	return &UserHttpController{
		registerUserCmdHdl: registerUserCmdHdl,
		authCmdHdl:         authCmdHdl,
		introspectCmdHdl:   introspectCmdHdl,
		generateCode:       generateCode,
		verifyCode:         verifyCode,
	}
}

func (ctrl *UserHttpController) SetupRoutes(g *gin.RouterGroup, authMld gin.HandlerFunc) {
	// Authentication group API
	g.POST("/register", ctrl.RegisterAPI)
	g.POST("/authenticate", ctrl.AuthenticateAPI) // Login
	g.GET("/profile", authMld, ctrl.GetProfileAPI)
	g.POST("/rpc/users/introspect-token", ctrl.IntrospectTokenRpcAPI) // RPC
	g.GET("/generateCode", authMld, ctrl.GenerateCodeAPI)
	g.GET("/verify/:code", authMld, ctrl.VerifyCodeAPI)

	// User info group API
	// users := g.Group("/users")
	// users.GET("", ctrl.List)
	// users.GET("/:id", ctrl.GetByID)
	// users.PATCH("/:id", ctrl.Update)
	// users.DELETE("/:id", ctrl.Delete)
}
