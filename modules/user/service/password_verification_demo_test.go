package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// TestPasswordVerificationDemo demonstrates that password verification is now working
func TestPasswordVerificationDemo(t *testing.T) {
	// Simulate the registration process
	plainPassword := "mySecurePassword123"
	salt := "randomsalt12345"
	
	// Hash password the same way as registration
	saltPass := fmt.Sprintf("%s.%s", salt, plainPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltPass), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// Create a user with the hashed password
	user := &usermodel.User{
		Id:       uuid.New(),
		Email:    "demo@example.com",
		Password: string(hashedPassword),
		Salt:     salt,
		Status:   datatype.StatusActive,
	}
	
	// Create mock repository with the user
	mockRepo := &mockAuthRepo{
		users: map[string]*usermodel.User{
			"demo@example.com": user,
		},
	}
	
	// Create mock token issuer
	mockTokenIssuer := &mockTokenIssuer{shouldFail: false}
	
	// Create authentication handler
	authHandler := NewAuthenticateCommandHandler(mockRepo, mockTokenIssuer)
	
	t.Run("correct password should authenticate", func(t *testing.T) {
		req := AuthenticateReq{
			Email:    "demo@example.com",
			Password: plainPassword, // Use the original plain password
		}
		
		result, err := authHandler.Execute(context.Background(), req)
		
		if err != nil {
			t.Errorf("Expected successful authentication, got error: %v", err)
			return
		}
		
		if result == nil {
			t.Error("Expected authentication result, got nil")
			return
		}
		
		if result.Token != "mock-jwt-token" {
			t.Errorf("Expected token 'mock-jwt-token', got '%s'", result.Token)
		}
		
		t.Logf("✅ Password verification successful! Token: %s", result.Token)
	})
	
	t.Run("wrong password should fail", func(t *testing.T) {
		req := AuthenticateReq{
			Email:    "demo@example.com",
			Password: "wrongPassword", // Wrong password
		}
		
		result, err := authHandler.Execute(context.Background(), req)
		
		if err == nil {
			t.Error("Expected authentication to fail with wrong password, but it succeeded")
			return
		}
		
		if result != nil {
			t.Error("Expected nil result for failed authentication")
			return
		}
		
		t.Logf("✅ Password verification correctly rejected wrong password: %v", err)
	})
	
	t.Run("empty password should fail", func(t *testing.T) {
		req := AuthenticateReq{
			Email:    "demo@example.com",
			Password: "", // Empty password
		}
		
		result, err := authHandler.Execute(context.Background(), req)
		
		if err == nil {
			t.Error("Expected authentication to fail with empty password, but it succeeded")
			return
		}
		
		if result != nil {
			t.Error("Expected nil result for failed authentication")
			return
		}
		
		t.Logf("✅ Password verification correctly rejected empty password: %v", err)
	})
	
	t.Run("password with different case should fail", func(t *testing.T) {
		req := AuthenticateReq{
			Email:    "demo@example.com",
			Password: "MYSECUREPASSWORD123", // Different case
		}
		
		result, err := authHandler.Execute(context.Background(), req)
		
		if err == nil {
			t.Error("Expected authentication to fail with different case password, but it succeeded")
			return
		}
		
		if result != nil {
			t.Error("Expected nil result for failed authentication")
			return
		}
		
		t.Logf("✅ Password verification is case-sensitive (correctly rejected): %v", err)
	})
}

// TestPasswordHashingConsistency verifies that our hashing is consistent with registration
func TestPasswordHashingConsistency(t *testing.T) {
	password := "testPassword123"
	salt := "testsalt"
	
	// Hash password the same way as in registration service
	saltPass := fmt.Sprintf("%s.%s", salt, password)
	hash1, err := bcrypt.GenerateFromPassword([]byte(saltPass), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}
	
	// Verify the password can be validated
	err = bcrypt.CompareHashAndPassword(hash1, []byte(saltPass))
	if err != nil {
		t.Errorf("Password verification failed: %v", err)
	}
	
	// Verify wrong password fails
	wrongSaltPass := fmt.Sprintf("%s.%s", salt, "wrongPassword")
	err = bcrypt.CompareHashAndPassword(hash1, []byte(wrongSaltPass))
	if err == nil {
		t.Error("Expected wrong password to fail verification")
	}
	
	t.Log("✅ Password hashing and verification is consistent")
}
