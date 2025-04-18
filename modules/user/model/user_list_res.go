package usermodel

import (
	"time"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UserListRes struct {
	Items      []UserSearchResDto    `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

type UserSearchResDto struct {
	Id        uuid.UUID  `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Role      UserRole   `json:"role"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
