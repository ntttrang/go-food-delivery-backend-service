package paymentmodule

import (
	"github.com/gin-gonic/gin"
	httpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/payment/infras/controller/http-gin"
	gormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/payment/infras/repository/gorm-mysql"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupPaymentModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup repositories
	cardRepo := gormmysql.NewCardRepo(dbCtx)

	// Setup card handlers
	createCardHandler := service.NewCreateCardCommandHandler(cardRepo)
	getCardByIDHandler := service.NewGetCardByIDQueryHandler(cardRepo)
	getCardsByUserIDHandler := service.NewGetCardsByUserIDQueryHandler(cardRepo)
	updateCardStatusHandler := service.NewUpdateCardStatusCommandHandler(cardRepo)

	// Setup controllers
	cardController := httpgin.NewCardController(
		createCardHandler,
		getCardByIDHandler,
		getCardsByUserIDHandler,
		updateCardStatusHandler,
	)

	// Setup routes
	cardController.SetupRoutes(g, appCtx.MiddlewareProvider().Auth())
}
