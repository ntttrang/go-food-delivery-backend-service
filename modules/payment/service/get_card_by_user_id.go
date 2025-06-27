package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
)

type IGetCardByUserIdRepo interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Card, error)
}

type GetCardsByUserIDQueryHandler struct {
	repo IGetCardByUserIdRepo
}

func NewGetCardsByUserIDQueryHandler(repo IGetCardByUserIdRepo) *GetCardsByUserIDQueryHandler {
	return &GetCardsByUserIDQueryHandler{repo: repo}
}

func (h *GetCardsByUserIDQueryHandler) Execute(ctx context.Context, userID uuid.UUID) ([]model.Card, error) {
	cards, err := h.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var activeCards []model.Card
	for _, card := range cards {
		if card.Status == string(model.CardStatusActive) {
			card.CardNumber = "**** **** **** " + card.CardNumber[len(card.CardNumber)-4:]
			card.CVV = "***"
			activeCards = append(activeCards, card)
		}
	}

	return activeCards, nil
}
