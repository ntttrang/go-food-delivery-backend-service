package service

// type CardCreateRequest struct {
// 	Method         string `json:"method"`
// 	Provider       string `json:"provider"`
// 	CardholderName string `json:"cardholderName"`
// 	CardNumber     string `json:"cardNumber"`
// 	CardType       string `json:"cardType"`
// 	CVV            string `json:"cvv"`
// 	ExpiryMonth    string `json:"expiryMonth"`
// 	ExpiryYear     string `json:"expiryYear"`

// 	UserID uuid.UUID `json:"_"` // from ctx
// }

// type CardCreateResponse struct {
// 	ID string `json:"id"`
// }

// type CardService struct {
// 	repo CardRepository
// }

// type CardRepository interface {
// 	Create(ctx context.Context, card *model.Card) error
// 	FindByID(ctx context.Context, id string) (*model.Card, error)
// 	UpdateStatusByID(ctx context.Context, id string, status string) error
// 	FindByUserID(ctx context.Context, userID string) ([]model.Card, error)
// }

// func NewCardService(repo CardRepository) *CardService {
// 	return &CardService{repo: repo}
// }

// func (s *CardService) CreateCard(ctx context.Context, req CardCreateRequest) (*CardCreateResponse, error) {
// 	id := uuid.NewString()
// 	card := &model.Card{
// 		ID:             id,
// 		Method:         req.Method,
// 		Provider:       req.Provider,
// 		CardholderName: req.CardholderName,
// 		CardNumber:     req.CardNumber,
// 		CardType:       req.CardType,
// 		CVV:            req.CVV,
// 		ExpiryMonth:    req.ExpiryMonth,
// 		ExpiryYear:     req.ExpiryYear,
// 		UserID:         req.UserID,
// 		Status:         "active",
// 	}
// 	err := s.repo.Create(ctx, card)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &CardCreateResponse{ID: id}, nil
// }

// func (s *CardService) GetCardByID(ctx context.Context, id string) (*model.Card, error) {
// 	return s.repo.FindByID(ctx, id)
// }

// func (s *CardService) UpdateCardStatusByID(ctx context.Context, id string, status string) error {
// 	return s.repo.UpdateStatusByID(ctx, id, status)
// }

// func (s *CardService) GetCardsByUserID(ctx context.Context, userID string) ([]model.Card, error) {
// 	return s.repo.FindByUserID(ctx, userID)
// }
