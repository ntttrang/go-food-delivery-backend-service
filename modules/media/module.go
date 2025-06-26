package mediamodule

import (
	"github.com/gin-gonic/gin"

	mediahttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/media/infras/controller/http-gin"
	gormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/media/infras/repository/gorm-mysql"
	mediaservice "github.com/ntttrang/go-food-delivery-backend-service/modules/media/service"
	sharedinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func InitializeMediaController(appCtx sharedinfras.IAppContext) *mediahttpgin.MediaHTTPController {
	dbCtx := appCtx.DbContext()

	// Setup repositories and services
	mediaRepository := gormmysql.NewImageRepository(dbCtx)

	// Set up command handlers
	createCommandHandler := mediaservice.NewCreateCommandHandler(mediaRepository)

	// Create HTTP controller
	mediaHTTPController := mediahttpgin.NewMediaHTTPController(createCommandHandler, appCtx.Uploader())
	return mediaHTTPController
}

func SetupMediaModule(appCtx sharedinfras.IAppContext, g *gin.RouterGroup) {
	mediaCtl := InitializeMediaController(appCtx)

	mediaCtl.SetupRoutes(g, appCtx.MiddlewareProvider().Auth())
}
