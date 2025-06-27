package model

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Payment method constants
const (
	MethodCreditCard = "CREDIT_CARD"
	MethodDebitCard  = "DEBIT_CARD"
)

// Payment provider constants
const (
	ProviderStripe = "STRIPE"
	ProviderPaypal = "PAYPAL"
)

// Card type constants
const (
	CardTypeVisa       = "VISA"
	CardTypeMastercard = "MASTERCARD"
	CardTypeJCB        = "JCB"
)

// Card status constants
type CardStatus string

const (
	CardStatusActive   CardStatus = "ACTIVE"
	CardStatusInactive CardStatus = "INACTIVE"
	CardStatusDeleted  CardStatus = "DELETED"
)

type Card struct {
	ID             uuid.UUID `gorm:"column:id" json:"id"`
	Method         string    `gorm:"column:method" json:"method"`
	Provider       string    `gorm:"column:provider" json:"provider"`
	CardholderName string    `gorm:"column:cardholder_name" json:"cardholderName"`
	CardNumber     string    `gorm:"column:card_number" json:"cardNumber"`
	CardType       string    `gorm:"column:card_type" json:"cardType"`
	ExpiryMonth    string    `gorm:"column:expiry_month" json:"expiryMonth"`
	ExpiryYear     string    `gorm:"column:expiry_year" json:"expiryYear"`
	CVV            string    `gorm:"column:cvv" json:"cvv"`
	UserID         uuid.UUID `gorm:"column:user_id" json:"userId"`
	Status         string    `gorm:"column:status" json:"status"`
	sharedmodel.DateDto
}

func (Card) TableName() string {
	return "cards"
}

// Validate validates the card data
func (c *Card) Validate() error {
	// Validate Method
	if c.Method == "" {
		return ErrMethodRequired
	}
	if c.Method != MethodCreditCard && c.Method != MethodDebitCard {
		return ErrInvalidPaymentMethod
	}

	// Validate Provider
	if c.Provider == "" {
		return ErrProviderRequired
	}
	if c.Provider != ProviderStripe && c.Provider != ProviderPaypal {
		return ErrInvalidPaymentMethod
	}

	// Validate CardholderName
	c.CardholderName = strings.TrimSpace(c.CardholderName)
	if c.CardholderName == "" {
		return ErrInvalidCardholderName
	}

	// Validate CardNumber (basic validation)
	c.CardNumber = strings.ReplaceAll(c.CardNumber, " ", "")
	if !regexp.MustCompile(`^\d{13,19}$`).MatchString(c.CardNumber) {
		return ErrInvalidCardNumber
	}

	// Validate CardType
	if c.CardType == "" {
		return ErrCardTypeRequired
	}
	if c.CardType != CardTypeVisa && c.CardType != CardTypeMastercard && c.CardType != CardTypeJCB {
		return ErrInvalidPaymentMethod
	}

	// Validate ExpiryMonth
	if c.ExpiryMonth == "" {
		return ErrInvalidExpiryDate
	}
	month, err := strconv.Atoi(c.ExpiryMonth)
	if err != nil || month < 1 || month > 12 {
		return ErrInvalidExpiryDate
	}

	// Validate ExpiryYear
	if c.ExpiryYear == "" {
		return ErrInvalidExpiryDate
	}
	year, err := strconv.Atoi(c.ExpiryYear)
	if err != nil {
		return ErrInvalidExpiryDate
	}

	// Check if card is expired
	currentYear := time.Now().Year()
	currentMonth := int(time.Now().Month())
	if year < currentYear || (year == currentYear && month < currentMonth) {
		return ErrInvalidExpiryDate
	}

	// Validate CVV
	if !regexp.MustCompile(`^\d{3,4}$`).MatchString(c.CVV) {
		return ErrInvalidCardCVV
	}

	// Validate UserID
	if c.UserID == uuid.Nil {
		return ErrUserIDRequired
	}

	return nil
}
