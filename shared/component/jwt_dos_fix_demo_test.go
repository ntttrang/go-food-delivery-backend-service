package sharecomponent

import (
	"context"
	"testing"
	"time"
)

// TestJWTDoSVulnerabilityFixed demonstrates that the JWT DoS vulnerability has been fixed
func TestJWTDoSVulnerabilityFixed(t *testing.T) {
	// Create JWT component
	jwtComp := NewJwtComp("test-secret-key", 3600)

	t.Run("invalid token should return error not crash app", func(t *testing.T) {
		// Test with completely invalid token
		invalidToken := "invalid.jwt.token"

		// This should return an error, not crash the application with log.Fatal
		userId, err := jwtComp.Validate(invalidToken)

		if err == nil {
			t.Error("Expected error for invalid token, got nil")
		}

		if userId != "" {
			t.Errorf("Expected empty userId for invalid token, got: %s", userId)
		}

		t.Logf("✅ Invalid token properly handled with error: %v", err)
	})

	t.Run("malformed token should return error not crash app", func(t *testing.T) {
		// Test with malformed token
		malformedToken := "malformed-token-without-dots"

		userId, err := jwtComp.Validate(malformedToken)

		if err == nil {
			t.Error("Expected error for malformed token, got nil")
		}

		if userId != "" {
			t.Errorf("Expected empty userId for malformed token, got: %s", userId)
		}

		t.Logf("✅ Malformed token properly handled with error: %v", err)
	})

	t.Run("expired token should return error not crash app", func(t *testing.T) {
		// Create a token that expires immediately
		shortLivedJwt := NewJwtComp("test-secret-key", -1) // Negative expiry = already expired

		// Generate a token
		token, err := shortLivedJwt.IssueToken(context.Background(), "test-user-id")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Wait a moment to ensure expiry
		time.Sleep(100 * time.Millisecond)

		// Try to validate the expired token
		userId, err := jwtComp.Validate(token)

		if err == nil {
			t.Error("Expected error for expired token, got nil")
		}

		if userId != "" {
			t.Errorf("Expected empty userId for expired token, got: %s", userId)
		}

		t.Logf("✅ Expired token properly handled with error: %v", err)
	})

	t.Run("token with wrong signature should return error not crash app", func(t *testing.T) {
		// Create token with one secret
		jwt1 := NewJwtComp("secret-1", 3600)
		token, err := jwt1.IssueToken(context.Background(), "test-user-id")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Try to validate with different secret
		jwt2 := NewJwtComp("secret-2", 3600)
		userId, err := jwt2.Validate(token)

		if err == nil {
			t.Error("Expected error for token with wrong signature, got nil")
		}

		if userId != "" {
			t.Errorf("Expected empty userId for token with wrong signature, got: %s", userId)
		}

		t.Logf("✅ Token with wrong signature properly handled with error: %v", err)
	})

	t.Run("empty token should return error not crash app", func(t *testing.T) {
		userId, err := jwtComp.Validate("")

		if err == nil {
			t.Error("Expected error for empty token, got nil")
		}

		if userId != "" {
			t.Errorf("Expected empty userId for empty token, got: %s", userId)
		}

		t.Logf("✅ Empty token properly handled with error: %v", err)
	})

	t.Run("valid token should work correctly", func(t *testing.T) {
		// Generate a valid token
		expectedUserId := "test-user-123"
		token, err := jwtComp.IssueToken(context.Background(), expectedUserId)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Validate the token
		userId, err := jwtComp.Validate(token)

		if err != nil {
			t.Errorf("Unexpected error for valid token: %v", err)
		}

		if userId != expectedUserId {
			t.Errorf("Expected userId %s, got %s", expectedUserId, userId)
		}

		t.Logf("✅ Valid token properly validated. UserId: %s", userId)
	})
}

// TestApplicationResilienceToConnectionFailures demonstrates that the app doesn't crash on connection failures
func TestApplicationResilienceToConnectionFailures(t *testing.T) {
	t.Run("NATS connection failure should not crash app", func(t *testing.T) {
		// This would previously cause log.Fatal and crash the app
		// Now it should return an error gracefully
		_, err := NewNatsComp()

		// We expect an error since NATS is likely not running in test environment
		if err == nil {
			t.Log("✅ NATS connection succeeded (NATS is running)")
		} else {
			t.Logf("✅ NATS connection failure handled gracefully: %v", err)
		}

		// The important thing is that we're still here and the test didn't crash
		t.Log("✅ Application continues to run despite NATS connection issues")
	})

	t.Run("MinIO connection failure should not crash app", func(t *testing.T) {
		// This would previously cause log.Fatalln and crash the app
		// Now it should return an error gracefully
		_, err := NewS3Uploader("invalid-key", "invalid-bucket", "invalid-domain", "invalid-region", "invalid-secret", false)

		// We expect an error since the credentials are invalid
		if err == nil {
			t.Error("Expected error for invalid MinIO credentials")
		} else {
			t.Logf("✅ MinIO connection failure handled gracefully: %v", err)
		}

		// The important thing is that we're still here and the test didn't crash
		t.Log("✅ Application continues to run despite MinIO connection issues")
	})
}

// TestGracefulDegradation demonstrates that services work even when dependencies are unavailable
func TestGracefulDegradation(t *testing.T) {
	t.Run("NATS publish should handle nil connection gracefully", func(t *testing.T) {
		// Create a nil NATS component (simulating connection failure)
		var natsComp *natsComp = nil

		// This should not panic or crash
		err := natsComp.Publish(nil, "test-topic", nil)

		if err != nil {
			t.Errorf("Expected nil error for graceful degradation, got: %v", err)
		}

		t.Log("✅ NATS publish handles nil connection gracefully")
	})
}
