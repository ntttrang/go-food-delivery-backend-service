package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/payment/service"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, req *service.CreateCardReq) (*service.CreateCardRes, error)
}

type IGetCardByIDQueryHandler interface {
	Execute(ctx context.Context, id uuid.UUID) (*model.Card, error)
}

type IUpdateStatusCommandHandler interface {
	Execute(ctx context.Context, req *service.CardUpdateStatusDto) error
}

type IGetByUserIdQueryHandler interface {
	Execute(ctx context.Context, userID uuid.UUID) ([]model.Card, error)
}

// CardController handles HTTP requests for cards
type CardController struct {
	createCardHandler       ICreateCommandHandler
	getCardByIDHandler      IGetCardByIDQueryHandler
	getCardsByUserIDHandler IGetByUserIdQueryHandler
	updateCardStatusHandler IUpdateStatusCommandHandler
}

// NewCardController creates a new card controller
func NewCardController(
	createCardHandler ICreateCommandHandler,
	getCardByIDHandler IGetCardByIDQueryHandler,
	getCardsByUserIDHandler IGetByUserIdQueryHandler,
	updateCardStatusHandler IUpdateStatusCommandHandler,
) *CardController {
	return &CardController{
		createCardHandler:       createCardHandler,
		getCardByIDHandler:      getCardByIDHandler,
		getCardsByUserIDHandler: getCardsByUserIDHandler,
		updateCardStatusHandler: updateCardStatusHandler,
	}
}

func (c *CardController) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	cards := router.Group("/cards")
	{
		cards.POST("", authMiddleware, c.CreateCard)
		cards.GET("/:id", c.GetCardByID)
		cards.PATCH("/:id", c.UpdateCardStatus)
		cards.GET("/user/:userId", c.GetCardsByUserID)
	}
}
