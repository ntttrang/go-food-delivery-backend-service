package usermodule

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	userHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/user/infras/controller/http-gin"
	userRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/user/infras/repository"
	userService "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"gorm.io/gorm"
)

func SetupUserModule(db *gorm.DB, g *gin.RouterGroup) {
	// Setup Controller
	// repo
	userRepo := userRepo.NewUserRepo(db)
	jwtComp := shareComponent.NewJwtComp(os.Getenv("JWT_SECRET_KEY"), 3600*24*7)
	// service
	registerCmdHdl := userService.NewRegisterUserCommandHandler(userRepo)
	authCmdHdl := userService.NewAuthenticateCommandHandler(userRepo, jwtComp)
	introspectCmdHdl := userService.NewIntrospectCommandHandler(jwtComp, userRepo)
	introspectCmdHdlWrapper := userService.NewIntrospectCmdHdlWrapper(introspectCmdHdl)
	// controller
	userCtrl := userHttpgin.NewUserHttpController(registerCmdHdl, authCmdHdl)

	// Setup router
	users := g.Group("/users")
	userCtrl.SetupRoutes(users, middleware.Auth(introspectCmdHdlWrapper))
}
