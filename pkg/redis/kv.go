package redis

import (
	"context"
	"time"
)

func (r *redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redis) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *redis) SetEx(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return r.client.Set(ctx, key, value, expire).Err()
}

func (r *redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redis) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return r.client.IncrBy(ctx, key, increment).Result()
}

func (r *redis) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *redis) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return r.client.DecrBy(ctx, key, decrement).Result()
}

func (r *redis) MGet(ctx context.Context, keys ...string) ([]*string, error) {
	if len(keys) == 0 {
		return []*string{}, nil
	}

	interfaceSli, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return []*string{}, err
	}

	strSlice := make([]*string, 0, len(interfaceSli))
	for _, v := range interfaceSli {
		value, ok := v.(string)
		if v == nil || !ok {
			strSlice = append(strSlice, nil)
		} else {
			strSlice = append(strSlice, &value)
		}
	}
	return strSlice, nil
}

func (r *redis) MSet(ctx context.Context, kvs map[string]interface{}) error {
	if len(kvs) == 0 {
		return nil
	}

	sli := make([]interface{}, 0, len(kvs)*2)
	for k, v := range kvs {
		sli = append(sli, k, v)
	}

	return r.client.MSet(ctx, sli...).Err()
}

func (r *redis) SetNx(ctx context.Context, key string, value interface{}) (bool, error) {
	return r.client.SetNX(ctx, key, value, time.Hour).Result()
}

func (r *redis) SetNxExpire(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}
