package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IGetCardRepo interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Card, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Card, error)
}

type GetCardByIDQueryHandler struct {
	repo IGetCardRepo
}

func NewGetCardByIDQueryHandler(repo IGetCardRepo) *GetCardByIDQueryHandler {
	return &GetCardByIDQueryHandler{repo: repo}
}

func (h *GetCardByIDQueryHandler) Execute(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	card, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	if card != nil && card.Status != string(model.CardStatusActive) {
		return nil, datatype.ErrNotFound.WithDebug("card not found")
	}

	card.CardNumber = "**** **** **** " + card.CardNumber[len(card.CardNumber)-4:]
	card.CVV = "***"
	return card, nil
}
