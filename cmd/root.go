package cmd

import (
	"context"
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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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

		serviceName := os.Getenv("SERVICE_NAME")

		r := gin.Default()

		// Middleware
		r.Use(middleware.Recover())
		r.Use(otelgin.Middleware(serviceName))

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

		// Initialize the tracer
		tp, err := initTracer(serviceName)
		if err != nil {
			log.Fatalf("Failed to initialize tracer: %v", err)
		}

		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
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

func initTracer(serviceName string) (*trace.TracerProvider, error) {
	log.Println("Initializing trace with OTLP gRPC exporter")

	// Create OTLP gRPC exporter for Jaeger
	// Note: deprecated Jaeger HTTP => Don't use port 4318
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "localhost:4317"
	}
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
		otlptracegrpc.WithInsecure(),
	)

	if err != nil {
		log.Printf("Failed to create OTLP gRPC exporter: %v \n ", err)
		return nil, err
	}
	log.Println("OTLP gRPC exporter created successfully")

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("development"),
		),
	)

	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return nil, err
	}
	log.Println("Resource created successfully")

	// Create trace provider with the exporter
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set global propagator for distributed tracing context
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Println("Tracer initialized successfully")
	return tp, nil
}
