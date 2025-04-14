package usermodule

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	userHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/user/infras/controller/http-gin"
	userRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/user/infras/repository"
	userService "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupUserModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup Controller
	// repo
	userRepo := userRepo.NewUserRepo(dbCtx)
	jwtComp := shareComponent.NewJwtComp(os.Getenv("JWT_SECRET_KEY"), 3600*24*7)
	// service
	registerCmdHdl := userService.NewRegisterUserCommandHandler(userRepo)
	authCmdHdl := userService.NewAuthenticateCommandHandler(userRepo, jwtComp)
	introspectCmdHdl := userService.NewIntrospectCommandHandler(jwtComp, userRepo)
	introspectCmdHdlWrapper := userService.NewIntrospectCmdHdlWrapper(introspectCmdHdl)

	redisCache := shareinfras.NewRedisAdapter(appCtx.GetConfig().RedisConfig)
	email := shareComponent.NewEmailService(appCtx.GetConfig().EmailConfig)
	generateCode := userService.NewGenerateCode(userRepo, redisCache, email)
	verifyCode := userService.NewVerifyCode(userRepo, redisCache)
	// controller
	userCtrl := userHttpgin.NewUserHttpController(registerCmdHdl, authCmdHdl, introspectCmdHdl, generateCode, verifyCode)

	// Setup router
	userCtrl.SetupRoutes(g, middleware.Auth(introspectCmdHdlWrapper))
}
