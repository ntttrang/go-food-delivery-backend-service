package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	usermodule "github.com/ntttrang/go-food-delivery-backend-service/modules/user"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
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

		v1 := r.Group("/v1")
		appCtx := shareinfras.NewAppContext(db)

		usermodule.SetupUserModule(appCtx, v1)

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
