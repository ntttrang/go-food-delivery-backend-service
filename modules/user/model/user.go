package usermodel

import (
	"github.com/google/uuid"
	datatype "github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type User struct {
	Id        uuid.UUID           `gorm:"id"`
	Email     string              `gorm:"email"`
	FbId      string              `gorm:"fb_id"`
	GgId      string              `gorm:"gg_id"`
	Password  string              `gorm:"password"`
	Salt      string              `gorm:"salt"`
	LastName  string              `gorm:"last_name"`
	FirstName string              `gorm:"first_name"`
	Phone     string              `gorm:"phone"`
	Role      datatype.UserRole   `gorm:"role"`
	Type      datatype.UserType   `gorm:"type"`
	Avatar    string              `gorm:"avatar"`
	Status    datatype.UserStatus `gorm:"status"`
	sharedmodel.DateDto
}

func (User) TableName() string {
	return "users"
}
