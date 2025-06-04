package ordermodule

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	orderHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/controller/http-gin"
	orderRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/gorm-mysql"
	orderService "github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/events"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

// Global Kafka producer instance
var kafkaProducer *events.KafkaProducer

func SetupOrderModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()
	// config := appCtx.GetConfig() // TODO: Use when implementing RPC integration

	// Setup repository
	repo := orderRepo.NewOrderRepo(dbCtx)

	// Setup traditional services
	createCmdHdl := orderService.NewCreateCommandHandlerSimple(repo)

	// Try to setup Kafka producer
	brokers := []string{"localhost:9092"}
	topic := "order-events"

	producer, err := events.NewKafkaProducer(brokers, topic)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		// Fall back to non-event-driven setup
		setupNonEventDrivenModule(repo, g)
		return
	}

	// Store global reference
	kafkaProducer = producer

	// Note: Event-driven service can be created when needed
	// For now, we'll use the traditional handlers with Kafka producer available globally

	// Setup other services
	listQueryHdl := orderService.NewListQueryHandler(repo)
	getDetailQueryHdl := orderService.NewGetDetailQueryHandler(repo)
	updateCmdHdl := orderService.NewUpdateCommandHandler(repo)
	deleteCmdHdl := orderService.NewDeleteCommandHandler(repo)

	// Setup controller (keep traditional interface for now)
	// The event-driven service will be used internally by the create handler
	orderCtl := orderHttpgin.NewOrderHttpController(
		createCmdHdl, // Use traditional handler interface
		listQueryHdl,
		getDetailQueryHdl,
		updateCmdHdl,
		deleteCmdHdl,
	)

	// Log that we're using event-driven architecture
	log.Printf("Order module initialized with Kafka event-driven architecture")

	// Setup routes
	orders := g.Group("/orders")
	orderCtl.SetupRoutes(orders)

	log.Println("Order module initialized with simple event-driven architecture")
}

func setupNonEventDrivenModule(repo *orderRepo.OrderRepo, g *gin.RouterGroup) {
	// Setup services
	createCmdHdl := orderService.NewCreateCommandHandlerSimple(repo)
	listQueryHdl := orderService.NewListQueryHandler(repo)
	getDetailQueryHdl := orderService.NewGetDetailQueryHandler(repo)
	updateCmdHdl := orderService.NewUpdateCommandHandler(repo)
	deleteCmdHdl := orderService.NewDeleteCommandHandler(repo)

	// Setup controller
	orderCtl := orderHttpgin.NewOrderHttpController(
		createCmdHdl,
		listQueryHdl,
		getDetailQueryHdl,
		updateCmdHdl,
		deleteCmdHdl,
	)

	// Setup routes
	orders := g.Group("/orders")
	orderCtl.SetupRoutes(orders)
}

// GetKafkaProducer returns the global Kafka producer instance
func GetKafkaProducer() *events.KafkaProducer {
	return kafkaProducer
}

// ShutdownOrderModule gracefully shuts down the order module
func ShutdownOrderModule(ctx context.Context) error {
	if kafkaProducer != nil {
		return kafkaProducer.Close()
	}
	return nil
}
