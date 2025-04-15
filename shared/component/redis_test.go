package shareComponent

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func TestRedisAdapter_Set(t *testing.T) {
	// Connect Redis
	config := datatype.GetConfig()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test redis connection",
			fields: fields{
				client: &redis.Client{},
			},
			args: args{
				key:        "test1",
				value:      "Hi, Trang 2",
				expiration: time.Second * 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewRedisAdapter(datatype.RedisConfig{
				Host:     config.RedisConfig.Host,
				Password: config.RedisConfig.Password,
			})
			if err := a.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("RedisAdapter.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
