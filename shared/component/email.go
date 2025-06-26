package sharecomponent

import (
	"fmt"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	cfg datatype.EmailConfig
}

func NewEmailService(cfg datatype.EmailConfig) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

func (e *EmailService) SendEmail(message sharedmodel.EmailMessage) error {
	config := e.cfg

	// Create a new message
	msg := gomail.NewMessage()

	// Set email headers
	msg.SetHeader("From", message.From)
	msg.SetHeader("To", message.To...)
	msg.SetHeader("Subject", message.Subject)

	msg.SetBody("text/html", message.Body)
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
