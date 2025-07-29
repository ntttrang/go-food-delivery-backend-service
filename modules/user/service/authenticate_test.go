package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Mock implementations for testing
type mockAuthRepo struct {
	users map[string]*usermodel.User
}

func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, usermodel.ErrUserNotFound
}

type mockTokenIssuer struct {
	shouldFail bool
}

func (m *mockTokenIssuer) IssueToken(ctx context.Context, userId string) (string, error) {
	if m.shouldFail {
		return "", errors.New("token generation failed")
	}
	return "mock-jwt-token", nil
}

func (m *mockTokenIssuer) ExpIn() int {
	return 3600
}

// Mock device token repository
type mockDeviceTokenRepo struct {
	shouldFail bool
}

func (m *mockDeviceTokenRepo) Insert(ctx context.Context, deviceToken *usermodel.UserDeviceToken) error {
	if m.shouldFail {
		return errors.New("device token insert failed")
	}
	return nil
}

func (m *mockDeviceTokenRepo) FindByToken(ctx context.Context, token string) (*usermodel.UserDeviceToken, error) {
	return nil, errors.New("not implemented")
}

func (m *mockDeviceTokenRepo) RevokeByUserId(ctx context.Context, userId string) error {
	return nil
}

func (m *mockDeviceTokenRepo) RevokeByToken(ctx context.Context, token string) error {
	return nil
}

// Mock refresh token generator
type mockRefreshTokenGenerator struct {
	shouldFail bool
}

func (m *mockRefreshTokenGenerator) GenerateRefreshToken() (string, error) {
	if m.shouldFail {
		return "", errors.New("refresh token generation failed")
	}
	return "mock-refresh-token", nil
}

func (m *mockRefreshTokenGenerator) RefreshTokenExpiry() time.Duration {
	return 30 * 24 * time.Hour // 30 days
}

func TestAuthenticateCommandHandler_Execute(t *testing.T) {
	// Setup test data with properly hashed passwords
	userId := uuid.New()

	// Create a properly hashed password for testing
	salt := "testsalt123456"
	password := "password123"
	saltPass := fmt.Sprintf("%s.%s", salt, password)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(saltPass), bcrypt.DefaultCost)

	activeUser := &usermodel.User{
		Id:       userId,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Salt:     salt,
		Status:   datatype.StatusActive,
	}

	deletedUser := &usermodel.User{
		Id:       uuid.New(),
		Email:    "deleted@example.com",
		Password: "hashedpassword",
		Status:   datatype.StatusDeleted,
	}

	bannedUser := &usermodel.User{
		Id:       uuid.New(),
		Email:    "banned@example.com",
		Password: "hashedpassword",
		Status:   datatype.StatusBanned,
	}

	tests := []struct {
		name                      string
		req                       AuthenticateReq
		mockRepo                  *mockAuthRepo
		mockDeviceTokenRepo       *mockDeviceTokenRepo
		mockTokenIssuer           *mockTokenIssuer
		mockRefreshTokenGenerator *mockRefreshTokenGenerator
		wantErr                   bool
		wantToken                 string
	}{
		{
			name: "successful authentication",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "password123", // This matches the password used to create the hash
			},
			mockRepo: &mockAuthRepo{
				users: map[string]*usermodel.User{
					"test@example.com": activeUser,
				},
			},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   false,
			wantToken:                 "mock-jwt-token",
		},
		{
			name: "wrong password",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockRepo: &mockAuthRepo{
				users: map[string]*usermodel.User{
					"test@example.com": activeUser,
				},
			},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "user not found",
			req: AuthenticateReq{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			mockRepo: &mockAuthRepo{
				users: map[string]*usermodel.User{},
			},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "deleted user",
			req: AuthenticateReq{
				Email:    "deleted@example.com",
				Password: "password123",
			},
			mockRepo: &mockAuthRepo{
				users: map[string]*usermodel.User{
					"deleted@example.com": deletedUser,
				},
			},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "banned user",
			req: AuthenticateReq{
				Email:    "banned@example.com",
				Password: "password123",
			},
			mockRepo: &mockAuthRepo{
				users: map[string]*usermodel.User{
					"banned@example.com": bannedUser,
				},
			},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "token generation failure",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockRepo: &mockAuthRepo{
				users: map[string]*usermodel.User{
					"test@example.com": activeUser,
				},
			},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: true},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "invalid email format",
			req: AuthenticateReq{
				Email:    "invalid-email",
				Password: "password123",
			},
			mockRepo:                  &mockAuthRepo{users: map[string]*usermodel.User{}},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "short password",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "123",
			},
			mockRepo:                  &mockAuthRepo{users: map[string]*usermodel.User{}},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "empty email",
			req: AuthenticateReq{
				Email:    "",
				Password: "password123",
			},
			mockRepo:                  &mockAuthRepo{users: map[string]*usermodel.User{}},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
		{
			name: "empty password",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "",
			},
			mockRepo:                  &mockAuthRepo{users: map[string]*usermodel.User{}},
			mockDeviceTokenRepo:       &mockDeviceTokenRepo{shouldFail: false},
			mockTokenIssuer:           &mockTokenIssuer{shouldFail: false},
			mockRefreshTokenGenerator: &mockRefreshTokenGenerator{shouldFail: false},
			wantErr:                   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewAuthenticateCommandHandler(tt.mockRepo, tt.mockDeviceTokenRepo, tt.mockTokenIssuer, tt.mockRefreshTokenGenerator)

			result, err := handler.Execute(context.Background(), tt.req, "Mozilla/5.0 (Test)")

			if tt.wantErr {
				if err == nil {
					t.Errorf("AuthenticateCommandHandler.Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("AuthenticateCommandHandler.Execute() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Errorf("AuthenticateCommandHandler.Execute() expected result, got nil")
				return
			}

			if result.AccessToken != tt.wantToken {
				t.Errorf("AuthenticateCommandHandler.Execute() token = %v, want %v", result.AccessToken, tt.wantToken)
			}

			if result.ExpIn != 3600 {
				t.Errorf("AuthenticateCommandHandler.Execute() expIn = %v, want %v", result.ExpIn, 3600)
			}
		})
	}
}

func TestAuthenticateReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     AuthenticateReq
		wantErr bool
	}{
		{
			name: "valid request",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			req: AuthenticateReq{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			req: AuthenticateReq{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "short password",
			req: AuthenticateReq{
				Email:    "test@example.com",
				Password: "123",
			},
			wantErr: true,
		},
		{
			name: "whitespace trimming",
			req: AuthenticateReq{
				Email:    "  test@example.com  ",
				Password: "  password123  ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("AuthenticateReq.Validate() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AuthenticateReq.Validate() unexpected error = %v", err)
				}

				// Check that whitespace was trimmed
				if tt.req.Email != "test@example.com" {
					t.Errorf("AuthenticateReq.Validate() email not trimmed properly: %v", tt.req.Email)
				}
				if tt.req.Password != "password123" {
					t.Errorf("AuthenticateReq.Validate() password not trimmed properly: %v", tt.req.Password)
				}
			}
		})
	}
}
