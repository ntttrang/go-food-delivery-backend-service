package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type FoodInsertDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`

	Id uuid.UUID `json:"-"`
}

func (c *FoodInsertDto) Validate() error {
	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return foodmodel.ErrNameRequired
	}

	return nil
}

func (c FoodInsertDto) ConvertToFood() *foodmodel.Food {
	return &foodmodel.Food{
		Name:        c.Name,
		Description: c.Description,
		Price:       c.Price,
	}
}

// Initilize service
type ICreateRepo interface {
	Insert(ctx context.Context, data *foodmodel.Food) error
}

type CreateCommandHandler struct {
	repo ICreateRepo
}

func NewCreateCommandHandler(repo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{repo: repo}
}

// Implement
func (s *CreateCommandHandler) Execute(ctx context.Context, data *FoodInsertDto) error {
	if err := data.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	food := data.ConvertToFood()
	food.Id, _ = uuid.NewV7()
	food.Status = sharedModel.StatusActive // Always set Active Status when insert

	if err := s.repo.Insert(ctx, food); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	data.Id = food.Id

	return nil
}
