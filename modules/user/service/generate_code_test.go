package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	sharemodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Mock implementations for testing
type mockUserRepo struct{}

func (m *mockUserRepo) FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error) {
	return &usermodel.User{
		Id:    id,
		Email: "test@example.com",
	}, nil
}

type mockRedisCache struct{}

func (m *mockRedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}

type mockEmail struct{}

func (m *mockEmail) SendEmail(message sharemodel.EmailMessage) error {
	return nil
}

func TestGenerateCode_Generate(t *testing.T) {
	userId := uuid.MustParse("019615db-9adb-7eff-ba03-45017274084c")

	// Create mock dependencies
	userRepo := &mockUserRepo{}
	redisCache := &mockRedisCache{}
	emailHdl := &mockEmail{}

	generateCode := NewGenerateCode(userRepo, redisCache, emailHdl)

	got, err := generateCode.Execute(context.Background(), userId)

	// Should not error
	if err != nil {
		t.Errorf("GenerateCode.Execute() error = %v, want nil", err)
		return
	}

	// Should return a 6-digit code
	if len(got) != 6 {
		t.Errorf("GenerateCode.Execute() returned code length = %d, want 6", len(got))
	}

	// Should be numeric
	for _, char := range got {
		if char < '0' || char > '9' {
			t.Errorf("GenerateCode.Execute() returned non-numeric code: %s", got)
			break
		}
	}
}
