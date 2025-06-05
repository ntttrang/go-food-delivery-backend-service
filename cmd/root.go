package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/gen/proto/category"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	cartmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/cart"
	categoryModule "github.com/ntttrang/go-food-delivery-backend-service/modules/category"
	categorygrpcctl "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/controller/grpc-ctrl"
	categorygormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/repository/gorm-mysql"
	foodmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/food"
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

// setupHTTPServer sets up and returns the HTTP server with all routes
func setupHTTPServer(appCtx shareinfras.IAppContext, port string) *http.Server {
	r := gin.Default()

	r.Use(middleware.Recover())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Static("/uploads", "./uploads")

	v1 := r.Group("/v1")

	// Setup all modules
	categoryModule.SetupCategoryModule(appCtx, v1)
	restaurantmodule.SetupRestaurantModule(appCtx, v1)
	usermodule.SetupUserModule(appCtx, v1)
	foodmodule.SetupFoodModule(appCtx, v1)
	mediamodule.SetupMediaModule(appCtx, v1)
	paymentmodule.SetupPaymentModule(appCtx, v1)
	cartmodule.SetupCartModule(appCtx, v1)
	ordermodule.SetupOrderModule(appCtx, v1)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
}

// setupGRPCServer sets up and returns the gRPC server
func setupGRPCServer(appCtx shareinfras.IAppContext) (*grpc.Server, net.Listener, error) {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Register GRPC services
	category.RegisterCategoryServer(s, categorygrpcctl.NewCategoryGrpcServer(categorygormmysql.NewCategoryRepo(appCtx.DbContext())))

	return s, lis, nil
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start service",
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}

		// Initialize application
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
			log.Printf("failed to connect database: %v \n", err)
		}
		log.Print("connected to database \n")

		appCtx := shareinfras.NewAppContext(db)

		// Setup HTTP server
		httpServer := setupHTTPServer(appCtx, port)

		// Setup gRPC server
		grpcServer, grpcListener, err := setupGRPCServer(appCtx)
		if err != nil {
			log.Fatalf("Failed to setup gRPC server: %v", err)
		}

		// Channel to listen for interrupt signal to terminate gracefully
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Use WaitGroup to wait for both servers to shutdown
		var wg sync.WaitGroup

		// Start gRPC server in a goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Starting gRPC server on :6000")
			if err := grpcServer.Serve(grpcListener); err != nil {
				log.Printf("gRPC server error: %v", err)
			}
		}()

		// Start HTTP server in a goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTP server on :%s", port)
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server error: %v", err)
			}
		}()

		// Wait for interrupt signal
		<-quit
		log.Println("Shutting down servers...")

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown HTTP server gracefully
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server forced to shutdown: %v", err)
		}

		// Shutdown gRPC server gracefully
		grpcServer.GracefulStop()

		// Wait for all servers to finish
		wg.Wait()
		log.Println("Servers stopped")
	},
}

func Execute() {
	// Add command
	// rootCmd.AddCommand(comsumerRestaurantCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("failed to execute command", err)
	}
}
