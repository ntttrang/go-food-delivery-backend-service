package usermodel

import (
	"time"

	"github.com/google/uuid"
)

type Verification struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;index"`
	Code      string    `gorm:"type:varchar(64);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp"`
}
