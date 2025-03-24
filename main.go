package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	categoryModule "github.com/ntttrang/go-food-delivery-backend-service/modules/category"
	restaurantmodule "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant"
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

	r.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")

	{
		categories := v1.Group("/categories")
		categoryModule.SetupCategoryModule(db, categories)
	}

	{
		restaurants := v1.Group("/restaurants")
		restaurantmodule.SetupRestaurantModule(db, restaurants)
	}

	r.Run(fmt.Sprintf(":%s", port))

}
