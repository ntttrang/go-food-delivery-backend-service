package usermodule

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	userHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/user/infras/controller/http-gin"
	repo "github.com/ntttrang/go-food-delivery-backend-service/modules/user/infras/repository/gorm-mysql"
	userService "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	sharecomponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupUserModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup Controller
	// repo
	userRepo := repo.NewUserRepo(dbCtx)
	userAddrRepo := repo.NewUserAddressRepo(dbCtx)
	jwtComp := sharecomponent.NewJwtComp(os.Getenv("JWT_SECRET_KEY"), 3600*24*7)
	ggOAuth := sharecomponent.NewGoogleOauth(appCtx.GetConfig().GoogleConfig)
	// service
	registerCmdHdl := userService.NewRegisterUserCommandHandler(userRepo)
	signUpGgCmdHdl := userService.NewSignUpGoogleCommandHandler(userRepo, jwtComp, ggOAuth)
	authCmdHdl := userService.NewAuthenticateCommandHandler(userRepo, jwtComp)
	introspectCmdHdl := userService.NewIntrospectCommandHandler(jwtComp, userRepo)
	introspectCmdHdlWrapper := userService.NewIntrospectCmdHdlWrapper(introspectCmdHdl)

	redisCache := sharecomponent.NewRedisAdapter(appCtx.GetConfig().RedisConfig)
	email := sharecomponent.NewEmailService(appCtx.GetConfig().EmailConfig)
	generateCode := userService.NewGenerateCode(userRepo, redisCache, email)
	verifyCode := userService.NewVerifyCode(userRepo, redisCache)

	listQueryHdl := userService.NewListQueryHandler(userRepo)
	getDetailQueryHdl := userService.NewGetDetailQueryHandler(userRepo)
	createCmdHdl := userService.NewCreateCommandHandler(userRepo)
	updateCmdHdl := userService.NewUpdateCommandHandler(userRepo)

	listAddrQueryHdl := userService.NewListAddrQueryHandler(userAddrRepo)
	createAddrCmdHdl := userService.NewCreateUserAddrCommandHandler(userAddrRepo)

	// controller
	userCtrl := userHttpgin.NewUserHttpController(
		registerCmdHdl, signUpGgCmdHdl, authCmdHdl, introspectCmdHdl,
		generateCode, verifyCode,
		listQueryHdl, getDetailQueryHdl, createCmdHdl, updateCmdHdl,
		userRepo, // user RPC
		listAddrQueryHdl, createAddrCmdHdl,
	)

	// Setup router
	userCtrl.SetupRoutes(g, middleware.Auth(introspectCmdHdlWrapper))
}
