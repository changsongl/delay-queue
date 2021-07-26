package redis

import (
	"context"
	gredis "github.com/go-redis/redis/v8"
	"time"
)

// IsError check if redis error, exclude redis nil
func IsError(err error) bool {
	return err != nil && err != gredis.Nil
}

// IsNil check if it is redis nil
func IsNil(err error) bool {
	return err == gredis.Nil
}

// Del a key
func (r *redis) Del(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Del(ctx, key).Result()
	return result == 1, err
}

// Exists check a key is exists
func (r *redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result == 1, err
}

// Expire a key
func (r *redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

// ExpireAt set a key expire time
func (r *redis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return r.client.ExpireAt(ctx, key, tm).Result()
}

// TTL get time to live
func (r *redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// FlushDB flush all data
func (r *redis) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}
