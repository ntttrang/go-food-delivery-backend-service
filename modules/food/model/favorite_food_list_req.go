package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FavoriteFoodListReq struct {
	UserId uuid.UUID `json:"-" form:"-"`
	sharedModel.PagingDto
}
