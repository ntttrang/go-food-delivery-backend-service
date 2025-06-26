package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type CategoryInsertDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Id        uuid.UUID          `json:"-"`
	Requester datatype.Requester `json:"-"`
}

func (c *CategoryInsertDto) Validate() error {
	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return categorymodel.ErrNameRequired
	}

	return nil
}

func (c CategoryInsertDto) ConvertToCategory() *categorymodel.Category {
	return &categorymodel.Category{
		Name:        c.Name,
		Description: c.Description,
	}
}

// Initilize service
type ICreateRepo interface {
	Insert(ctx context.Context, data *categorymodel.Category) error
}

type CreateCommandHandler struct {
	catRepo ICreateRepo
}

func NewCreateCommandHandler(catRepo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{catRepo: catRepo}
}

// Implement
func (s *CreateCommandHandler) Execute(ctx context.Context, data *CategoryInsertDto) error {
	if err := data.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Authorization check
	if data.Requester == nil {
		return datatype.ErrUnauthorized.WithDebug(categorymodel.ErrRequesterRequired.Error())
	}

	// Check if requester is admin or user (not shipper)
	// Use type assertion to get the role as a string
	role := data.Requester.GetRole()

	// Check role
	isAuthorized := role == string(datatype.RoleAdmin) || role == string(datatype.RoleUser)
	if !isAuthorized {
		return datatype.ErrForbidden.WithDebug(categorymodel.ErrPermission.Error())
	}

	category := data.ConvertToCategory()
	category.Id, _ = uuid.NewV7()
	category.Status = string(datatype.StatusActive) // Always set Active Status when insert

	if err := s.catRepo.Insert(ctx, category); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	data.Id = category.Id

	return nil
}
