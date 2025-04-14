package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IGetUserRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type IGetRedisCache interface {
	Get(ctx context.Context, key string, dest interface{}) error
}

type VerifyCode struct {
	userRepo   IGetUserRepo
	redisCache IGetRedisCache
}

func NewVerifyCode(userRepo IGetUserRepo, redisCache IGetRedisCache) *VerifyCode {
	return &VerifyCode{
		userRepo:   userRepo,
		redisCache: redisCache,
	}
}

func (v *VerifyCode) Execute(ctx context.Context, userId uuid.UUID, code string) (bool, error) {
	user, err := v.userRepo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return false, datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}
		return false, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var emailAddr string
	if user != nil {
		emailAddr = user.Email
	}
	// Store in redis
	var codeInCache *string
	err = v.redisCache.Get(ctx, emailAddr, &codeInCache)
	if err != nil {
		return false, err
	}

	if codeInCache == nil || *codeInCache != code {
		return false, nil
	}

	return true, nil
}
