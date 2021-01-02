package redis

import (
	"context"
	gredis "github.com/go-redis/redis/v8"
	"time"
)

const (
	KeepTTL time.Duration = -1
)

func (r *redis) IsError(err error) bool {
	return err != nil && err != gredis.Nil
}

func (r *redis) IsNil(err error) bool {
	return err == gredis.Nil
}

func (r *redis) Del(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Del(ctx, key).Result()
	return result == 1, err
}

func (r *redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	return result == 1, err
}

func (r *redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.Client.Expire(ctx, key, expiration).Result()
}

func (r *redis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return r.Client.ExpireAt(ctx, key, tm).Result()
}

func (r *redis) Ttl(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}
