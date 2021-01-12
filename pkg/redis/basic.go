package redis

import (
	"context"
	gredis "github.com/go-redis/redis/v8"
	"time"
)

const (
	KeepTTL time.Duration = -1
)

func IsError(err error) bool {
	return err != nil && err != gredis.Nil
}

func IsNil(err error) bool {
	return err == gredis.Nil
}

func (r *redis) Del(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Del(ctx, key).Result()
	return result == 1, err
}

func (r *redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result == 1, err
}

func (r *redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

func (r *redis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return r.client.ExpireAt(ctx, key, tm).Result()
}

func (r *redis) Ttl(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}
