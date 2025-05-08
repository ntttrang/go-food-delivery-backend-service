package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
)

type IRegisterUserCommandHandler interface {
	Execute(ctx context.Context, req *service.RegisterUserReq) error
}

type IAuthenticateCommandHandler interface {
	Execute(ctx context.Context, req service.AuthenticateReq) (*service.AuthenticateRes, error)
}
type IntrospectCommandHandler interface {
	Execute(ctx context.Context, req service.IntrospectReq) (*service.IntrospectRes, error)
}

type IGenerateCode interface {
	Execute(ctx context.Context, userId uuid.UUID) (string, error)
}

type IVerifyCode interface {
	Execute(ctx context.Context, userId uuid.UUID, code string) (bool, error)
}

type IListQueryHandler interface {
	Execute(ctx context.Context, req service.UserListReq) (service.UserListRes, error)
}

type IGetDetailQueryHandler interface {
	Execute(ctx context.Context, req service.UserDetailReq) (service.UserSearchResDto, error)
}

type ICreateCommandHandler interface {
	Execute(ctx context.Context, req *service.CreateUserReq) error
}

type IUpdateCommandHandler interface {
	Execute(ctx context.Context, req service.UpdateUserReq) error
}

type ISignUpGoogleCommandHandler interface {
	GetAuthCodeUrl(ctx context.Context) string
	AuthenticateByGoogle(ctx context.Context, state string, code string) (*service.AuthenticateRes, error)
}

type IRepoRPCUser interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]usermodel.User, error)
}

type IListAddrQueryHandler interface {
	Execute(ctx context.Context, req service.UserAddrListReq) (service.UserAddrListRes, error)
}

type ICreateAddrCommandHandler interface {
	Execute(ctx context.Context, req *service.CreateUserAddrReq) error
}

type UserHttpController struct {
	registerUserCmdHdl IRegisterUserCommandHandler
	signUpGgCmdHdl     ISignUpGoogleCommandHandler
	authCmdHdl         IAuthenticateCommandHandler
	introspectCmdHdl   IntrospectCommandHandler
	generateCode       IGenerateCode
	verifyCode         IVerifyCode

	listQueryHdl      IListQueryHandler
	getDetailQueryHdl IGetDetailQueryHandler
	createCmdHdl      ICreateCommandHandler
	updateCmdHdl      IUpdateCommandHandler

	rpcUser IRepoRPCUser

	listAddrQueryHdl IListAddrQueryHandler
	createAddrCmdHdl ICreateAddrCommandHandler
}

func NewUserHttpController(registerUserCmdHdl IRegisterUserCommandHandler, signUpGgCmdHdl ISignUpGoogleCommandHandler, authCmdHdl IAuthenticateCommandHandler, introspectCmdHdl IntrospectCommandHandler,
	generateCode IGenerateCode, verifyCode IVerifyCode,
	listQueryHdl IListQueryHandler, getDetailQueryHdl IGetDetailQueryHandler, createCmdHdl ICreateCommandHandler, updateCmdHdl IUpdateCommandHandler,
	rpcUser IRepoRPCUser,
	listAddrQueryHdl IListAddrQueryHandler, createAddrCmdHdl ICreateAddrCommandHandler) *UserHttpController {
	return &UserHttpController{
		registerUserCmdHdl: registerUserCmdHdl,
		signUpGgCmdHdl:     signUpGgCmdHdl,
		authCmdHdl:         authCmdHdl,
		introspectCmdHdl:   introspectCmdHdl,
		generateCode:       generateCode,
		verifyCode:         verifyCode,
		listQueryHdl:       listQueryHdl,
		getDetailQueryHdl:  getDetailQueryHdl,
		createCmdHdl:       createCmdHdl,
		updateCmdHdl:       updateCmdHdl,
		rpcUser:            rpcUser,
		listAddrQueryHdl:   listAddrQueryHdl,
		createAddrCmdHdl:   createAddrCmdHdl,
	}
}

func (ctrl *UserHttpController) SetupRoutes(g *gin.RouterGroup, authMld gin.HandlerFunc) {
	// Signup by email
	g.POST("/register", ctrl.RegisterAPI)
	// Sign Up with Google
	g.POST("/google/signup", ctrl.SignUpWithGoogleAPI)
	g.GET("/google/callback", ctrl.CallbackAPI)

	g.POST("/authenticate", ctrl.AuthenticateAPI) // Login
	g.GET("/profile", authMld, ctrl.GetProfileAPI)
	g.POST("/rpc/users/introspect-token", ctrl.IntrospectTokenRpcAPI) // RPC
	g.GET("/generate-code", authMld, ctrl.GenerateCodeAPI)
	g.GET("/verify/:code", authMld, ctrl.VerifyCodeAPI)
	//g.GET("/reset-password", ctrl.VerifyCodeAPI)

	// RPC
	g.POST("/rpc/users/find-by-ids", ctrl.RPCGetByIds)

	// User info group API
	users := g.Group("/users")
	users.POST("", ctrl.CreateUserAPI)
	users.GET("", ctrl.ListUsersAPI)
	users.GET("/:id", ctrl.GetUserDetailAPI)
	users.PATCH("/:id", ctrl.UpdateUseAPI)

	// Address
	users.POST("/address", authMld, ctrl.CreateUserAddrAPI)
	users.GET("/address", ctrl.ListUserAddrAPI)
	//users.GET("/address/:id", ctrl.GetUserAddrDetailAPI)
	//users.PATCH("/address/:id", ctrl.UpdateUserAddrAPI)
	//users.DELETE("/address/:id", ctrl.DeleteAddrAPI)

}
