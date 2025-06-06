package cmd

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	orderRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/gorm-mysql"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/rpc-client"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var consumerChangePaymentStatusCmd = &cobra.Command{

	Use:   "change-payment-status-cmd",
	Short: "Start consumer send email when change payment status",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := os.Getenv("DB_DSN")
		dbMaster, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Fatal("failed to connect database", err)
		}

		db := dbMaster.Debug()

		nc, err := nats.Connect(os.Getenv("NATS_URL"))

		if err != nil {
			log.Fatal("failed to connect nats", err)
		}

		appCtx := shareinfras.NewAppContext(db)
		dbCtx := appCtx.DbContext()
		orderRepo := orderRepo.NewOrderRepo(dbCtx)
		restaurantRpcClientRepo := rpcclient.NewRestaurantRPCClient(appCtx.GetConfig().RestaurantServiceURL)
		userRpcClientRepo := rpcclient.NewUserRPCClient(appCtx.GetConfig().UserServiceURL)
		emailSvc := shareComponent.NewEmailService(appCtx.GetConfig().EmailConfig)

		notificationService := service.NewOrderNotificationService(
			orderRepo,
			userRpcClientRepo,
			restaurantRpcClientRepo,
			emailSvc,
		)

		nc.Subscribe(datatype.EvtNotifyPaymentStatusChange, func(msg *nats.Msg) {
			log.Println("Subscribe: CHANGE PAYMENT STATUS")
			type evtNotifyChangePaymentStatusMsg struct {
				OrderID       string
				PaymentStatus string
			}

			var data evtNotifyChangePaymentStatusMsg

			if err := json.Unmarshal(msg.Data, &data); err != nil {
				log.Println("failed to unmarshal data", err)
				return
			}

			notificationService.NotifyPaymentStatusChange(context.Background(), data.OrderID, data.PaymentStatus)

			log.Printf("Send email notification to parties: %v \n", data)
		})

		// Setup graceful shutdown
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		// Block until we receive a signal
		log.Println("Consumer started. Press Ctrl+C to exit...")
		<-signalChan

		log.Println("Shutting down consumer...")

		// Drain connection (process pending messages before closing)
		if err := nc.Drain(); err != nil {
			log.Printf("Error draining NATS connection: %v", err)
		}

		// Close NATS connection
		nc.Close()

		log.Println("Consumer shutdown complete")
	},
}
