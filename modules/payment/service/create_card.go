package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type CreateCardReq struct {
	Method         string    `json:"method"`
	Provider       string    `json:"provider"`
	CardholderName string    `json:"cardholderName"`
	CardNumber     string    `json:"cardNumber"`
	CardType       string    `json:"cardType"`
	Cvv            string    `json:"cvv"`
	ExpiryMonth    string    `json:"expiryMonth"`
	ExpiryYear     string    `json:"expiryYear"`
	UserID         uuid.UUID `json:"-"`
}

type CreateCardRes struct {
	ID uuid.UUID `json:"id"`
}

type ICreateCardRepo interface {
	Create(ctx context.Context, card *model.Card) error
}

type CreateCardCommandHandler struct {
	repo ICreateCardRepo
}

func NewCreateCardCommandHandler(repo ICreateCardRepo) *CreateCardCommandHandler {
	return &CreateCardCommandHandler{repo: repo}
}

func (h *CreateCardCommandHandler) Execute(ctx context.Context, req *CreateCardReq) (*CreateCardRes, error) {
	card := &model.Card{
		Method:         req.Method,
		Provider:       req.Provider,
		CardholderName: strings.TrimSpace(req.CardholderName),
		CardNumber:     req.CardNumber,
		CardType:       req.CardType,
		CVV:            req.Cvv,
		ExpiryMonth:    req.ExpiryMonth,
		ExpiryYear:     req.ExpiryYear,
		UserID:         req.UserID,
		Status:         string(model.CardStatusActive),
	}

	// Validate the card
	if err := card.Validate(); err != nil {
		return nil, err
	}

	// Generate a new UUID for the card
	card.ID, _ = uuid.NewV7()

	// Create the card in the repository
	if err := h.repo.Create(ctx, card); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Return the response
	return &CreateCardRes{ID: card.ID}, nil
}
