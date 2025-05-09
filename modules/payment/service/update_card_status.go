package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type CardUpdateStatusDto struct {
	ID     uuid.UUID `json:"-"`
	Status string    `json:"status"`
}

type IUpdateCardRepo interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Card, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type UpdateCardStatusCommandHandler struct {
	repo IUpdateCardRepo
}

func NewUpdateCardStatusCommandHandler(repo IUpdateCardRepo) *UpdateCardStatusCommandHandler {
	return &UpdateCardStatusCommandHandler{repo: repo}
}

// Execute executes the update card status command
func (h *UpdateCardStatusCommandHandler) Execute(ctx context.Context, req *CardUpdateStatusDto) error {
	// Validate the status
	if req.Status != model.CardStatusActive && req.Status != model.CardStatusInactive && req.Status != model.CardStatusDeleted {
		return datatype.ErrBadRequest.WithDebug("invalid card status")
	}

	// Check if the card exists
	_, err := h.repo.FindByID(ctx, req.ID)
	if err != nil {
		return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
	}

	// Update the card status
	if err := h.repo.UpdateStatus(ctx, req.ID, req.Status); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
