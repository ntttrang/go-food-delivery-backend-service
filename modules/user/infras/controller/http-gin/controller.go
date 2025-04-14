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

type IListQueryHandler interface {
	Execute(ctx context.Context, req usermodel.UserListReq) (usermodel.UserListRes, error)
}

type IGetDetailQueryHandler interface {
	Execute(ctx context.Context, req usermodel.UserDetailReq) (usermodel.UserSearchResDto, error)
}

type ICreateCommandHandler interface {
	Execute(ctx context.Context, req *usermodel.CreateUserReq) error
}

type IUpdateCommandHandler interface {
	Execute(ctx context.Context, req usermodel.UpdateUserReq) error
}

type UserHttpController struct {
	registerUserCmdHdl IRegisterUserCommandHandler
	authCmdHdl         IAuthenticateCommandHandler
	introspectCmdHdl   IntrospectCommandHandler
	generateCode       IGenerateCode
	verifyCode         IVerifyCode
	listQueryHdl       IListQueryHandler
	getDetailQueryHdl  IGetDetailQueryHandler
	createCmdHdl       ICreateCommandHandler
	updateCmdHdl       IUpdateCommandHandler
}

func NewUserHttpController(registerUserCmdHdl IRegisterUserCommandHandler, authCmdHdl IAuthenticateCommandHandler, introspectCmdHdl IntrospectCommandHandler,
	generateCode IGenerateCode, verifyCode IVerifyCode,
	listQueryHdl IListQueryHandler, getDetailQueryHdl IGetDetailQueryHandler, createCmdHdl ICreateCommandHandler, updateCmdHdl IUpdateCommandHandler) *UserHttpController {
	return &UserHttpController{
		registerUserCmdHdl: registerUserCmdHdl,
		authCmdHdl:         authCmdHdl,
		introspectCmdHdl:   introspectCmdHdl,
		generateCode:       generateCode,
		verifyCode:         verifyCode,
		listQueryHdl:       listQueryHdl,
		getDetailQueryHdl:  getDetailQueryHdl,
		createCmdHdl:       createCmdHdl,
		updateCmdHdl:       updateCmdHdl,
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

	users := g.Group("/users")
	users.POST("", ctrl.CreateUserAPI)
	users.GET("", ctrl.ListUsersAPI)
	users.GET("/:id", ctrl.GetUserDetailAPI)
	users.PATCH("/:id", ctrl.UpdateUseAPI)

}
