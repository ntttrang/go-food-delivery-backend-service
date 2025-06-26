package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start uservice",
	Run: func(cmd *cobra.Command, args []string) {
		// Get env
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		dsn := os.Getenv("DB_DSN")
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		log.Print("connected to database \n")

		// Start Gin Engine
		r := gin.Default()

		// Set Middleware
		r.Use(middleware.Recover())

		// API Route
		r.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK,
				gin.H{"message": "Ok"},
			)
		})

		// TODO
		fmt.Println(db)

		// Run app
		r.Run(fmt.Sprintf(":%s", port))

	},
}

func Execute() {
	// Add command if have

	// Start server
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("failed to execute command", err)
	}
}
