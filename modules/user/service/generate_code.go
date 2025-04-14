package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IUserRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type IRedisCache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type IEmail interface {
	SendEmail(message sharedModel.EmailMessage) error
}

type GenerateCode struct {
	userRepo   IUserRepo
	redisCache IRedisCache
	emailHdl   IEmail
}

func NewGenerateCode(userRepo IUserRepo, redisCache IRedisCache, emailHdl IEmail) *GenerateCode {
	return &GenerateCode{
		userRepo:   userRepo,
		redisCache: redisCache,
		emailHdl:   emailHdl,
	}
}

func (g *GenerateCode) Execute(ctx context.Context, userId uuid.UUID) (string, error) {

	user, err := g.userRepo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return "", datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}
		return "", datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var emailAddr string
	if user != nil {
		emailAddr = user.Email
	}

	// Generate code
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	random := rand.New(source)
	verifyCode := fmt.Sprintf("%06d", random.Intn(1000000)) // Generates a 6-digit code

	// Store in redis
	err = g.redisCache.Set(ctx, emailAddr, verifyCode, time.Hour*1)
	if err != nil {
		return "", err
	}

	// Send email
	var msg sharedModel.EmailMessage
	msg.From = "sender@gmail.com"
	msg.To = []string{"minhtrang.2106@gmail.com"}
	msg.Subject = "[FD-Testing-v2] Your Verification Code"
	msg.Body = fmt.Sprintf("Your verification code is: %s", verifyCode)
	err = g.emailHdl.SendEmail(msg)
	if err != nil {
		return "", err
	}

	return verifyCode, nil
}
