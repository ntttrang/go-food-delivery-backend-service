package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/category"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/food"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	cartmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/cart"
	categoryModule "github.com/ntttrang/go-food-delivery-backend-service/modules/category"
	categorygrpcctl "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/controller/grpc-ctrl"
	categorygormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/repository/gorm-mysql"
	foodmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/food"
	foodgrpcctl "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/controller/grpc-ctrl"
	foodgormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/gorm-mysql"
	foodservice "github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	mediamodule "github.com/ntttrang/go-food-delivery-backend-service/modules/media"
	ordermodule "github.com/ntttrang/go-food-delivery-backend-service/modules/order"
	paymentmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/payment"
	restaurantmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant"
	usermodule "github.com/ntttrang/go-food-delivery-backend-service/modules/user"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start service",
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
		dsn := os.Getenv("DB_DSN")
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		log.Print("connected to database \n")

		r := gin.Default()

		r.Use(middleware.Recover())

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		r.Static("/uploads", "./uploads")

		v1 := r.Group("/v1")
		appCtx := shareinfras.NewAppContext(db)

		categoryModule.SetupCategoryModule(appCtx, v1)
		restaurantmodule.SetupRestaurantModule(appCtx, v1)
		usermodule.SetupUserModule(appCtx, v1)
		foodmodule.SetupFoodModule(appCtx, v1)
		mediamodule.SetupMediaModule(appCtx, v1)
		paymentmodule.SetupPaymentModule(appCtx, v1)
		cartmodule.SetupCartModule(appCtx, v1)
		ordermodule.SetupOrderModule(appCtx, v1)

		// Run gRPC server
		go func() {

			grpcPort := os.Getenv("GRPC_PORT")
			if grpcPort == "" {
				grpcPort = "6000"
			}

			// Create a listener on TCP port
			lis, err := net.Listen("tcp", ":"+grpcPort)
			if err != nil {
				log.Fatalln("Failed to listen:", err)
			}

			// Create a gRPC server object
			s := grpc.NewServer()

			// Register GRPC
			category.RegisterCategoryServer(s, categorygrpcctl.NewCategoryGrpcServer(categorygormmysql.NewCategoryRepo(appCtx.DbContext())))

			// Setup food gRPC server with update service
			foodRepo := foodgormmysql.NewFoodRepo(appCtx.DbContext())
			updateService := foodservice.NewUpdateCommandHandler(foodRepo)
			food.RegisterFoodServer(s, foodgrpcctl.NewFoodGrpcServer(foodRepo, updateService))

			// Serve gRPC Server
			log.Printf("Serving gRPC on 0.0.0.0:%s \n", grpcPort)
			log.Fatal(s.Serve(lis))
		}()

		// Run app server
		r.Run(fmt.Sprintf(":%s", port))

	},
}

func Execute() {
	setupConsumerCmd()
	// Add command
	rootCmd.AddCommand(consumerCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("failed to execute command", err)
	}
}
