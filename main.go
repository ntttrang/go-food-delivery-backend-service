package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	categoryModule "github.com/ntttrang/go-food-delivery-backend-service/modules/category"
	foodmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/food"
	restaurantmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant"
	usermodule "github.com/ntttrang/go-food-delivery-backend-service/modules/user"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
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

	v1 := r.Group("/v1")
	appCtx := shareinfras.NewAppContext(db)

	categoryModule.SetupCategoryModule(appCtx, v1)
	restaurantmodule.SetupRestaurantModule(appCtx, v1)
	usermodule.SetupUserModule(appCtx, v1)
	foodmodule.SetupFoodModule(appCtx, v1)

	r.Run(fmt.Sprintf(":%s", port))

}
