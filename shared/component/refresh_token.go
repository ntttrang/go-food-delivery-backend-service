package sharecomponent

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type RefreshTokenGenerator struct {
	expiry time.Duration
}

func NewRefreshTokenGenerator(expiry time.Duration) *RefreshTokenGenerator {
	return &RefreshTokenGenerator{
		expiry: expiry,
	}
}

// GenerateRefreshToken generates a cryptographically secure random refresh token
func (r *RefreshTokenGenerator) GenerateRefreshToken() (string, error) {
	// Generate 32 bytes of random data (256 bits of entropy)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64 raw URL-safe string (no padding, URL-safe)
	// This produces a 43-character string that's safe for URLs and JSON
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// RefreshTokenExpiry returns the refresh token expiry duration
func (r *RefreshTokenGenerator) RefreshTokenExpiry() time.Duration {
	return r.expiry
}
