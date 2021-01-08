package redis

import "context"

func (r *redis) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.LPush(ctx, key, values...).Result()
}

func (r *redis) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.RPush(ctx, key, values...).Result()
}

func (r *redis) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

func (r *redis) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

func (r *redis) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

func (r *redis) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

func (r *redis) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return r.client.LRem(ctx, key, count, value).Result()
}

func (r *redis) LIndex(ctx context.Context, key string, idx int64) (string, error) {
	return r.client.LIndex(ctx, key, idx).Result()
}

func (r *redis) LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	return r.client.LTrim(ctx, key, start, stop).Result()
}
