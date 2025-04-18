package shareComponent

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type JwtComp struct {
	secretKey string
	expIn     int
}

func NewJwtComp(secretKey string, expIn int) *JwtComp {
	return &JwtComp{
		secretKey: secretKey,
		expIn:     expIn,
	}
}

func (j *JwtComp) IssueToken(ctx context.Context, userId string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.expIn))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return tokenString, nil
}

func (j *JwtComp) ExpIn() int {
	return j.expIn
}

func (j *JwtComp) Validate(tokenStr string) (string, error) {
	var rc jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return rc.Subject, nil
}
