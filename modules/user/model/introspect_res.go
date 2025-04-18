package usermodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IntrospectRes struct {
	Id        uuid.UUID  `json:"id"`
	LastName  string     `json:"last_name"`
	FirstName string     `json:"first_name"`
	Role      UserRole   `json:"role"`
	Type      UserType   `json:"type"`
	Status    UserStatus `json:"status"`
	sharedModel.DateDto
}

func (ir *IntrospectRes) GetRole() uuid.UUID {
	return uuid.MustParse(string(ir.Role))
}

func (ir *IntrospectRes) Subject() uuid.UUID {
	return ir.Id
}
