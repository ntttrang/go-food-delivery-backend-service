package usermodel

import (
	"encoding/json"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UserType string
type UserRole string
type UserStatus string

const (
	TypeEmailPassword UserType   = "EMAIL_PASSWORD"
	TypeFacebook      UserType   = "FACEBOOK"
	TypeGmail         UserType   = "GMAIL"
	RoleUser          UserRole   = "USER"
	RoleAdmin         UserRole   = "ADMIN"
	RoleShipper       UserRole   = "SHIPPER"
	StatusPending     UserStatus = "PENDING"
	StatusActive      UserStatus = "ACTIVE"
	StatusInactive    UserStatus = "INACTIVE"
	StatusBanned      UserStatus = "BANNED"
	StatusDeleted     UserStatus = "DELETED"
)

type User struct {
	Id        uuid.UUID       `gorm:"id"`
	Email     string          `gorm:"email"`
	FbId      string          `gorm:"fb_id"`
	GgId      string          `gorm:"gg_id"`
	Password  string          `gorm:"password"`
	Salt      string          `gorm:"salt"`
	LastName  string          `gorm:"last_name"`
	FirstName string          `gorm:"first_name"`
	Phone     string          `gorm:"phone"`
	Role      UserRole        `gorm:"role"`
	Type      UserType        `gorm:"type"`
	Avatar    json.RawMessage `gorm:"avatar"`
	Status    UserStatus      `gorm:"status"`
	sharedModel.DateDto
}

func (User) TableName() string {
	return "users"
}
