package shareinfras

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(cfg datatype.RedisConfig) *RedisAdapter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       0, // default
	})

	return &RedisAdapter{client: rdb}
}

var ErrCacheMiss = errors.New("cache miss")

func (a *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return a.client.Set(key, data, expiration).Err()
}

func (a *RedisAdapter) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := a.client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (a *RedisAdapter) Delete(ctx context.Context, key string) error {
	return a.client.Del(key).Err()
}
