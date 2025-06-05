package ordermodel

import (
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Role      datatype.UserRole
}
