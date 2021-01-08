package redis

import (
	"context"
)

func (r *redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *redis) HMGet(ctx context.Context, key string, fields []string) ([]*string, error) {
	if len(fields) == 0 {
		return []*string{}, nil
	}

	objSli, err := r.client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return []*string{}, err
	}

	strSlice := make([]*string, 0, len(objSli))
	for _, v := range objSli {
		value, ok := v.(string)
		if v == nil || !ok {
			strSlice = append(strSlice, nil)
		} else {
			strSlice = append(strSlice, &value)
		}
	}
	return strSlice, nil
}

func (r *redis) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *redis) HMSet(ctx context.Context, key string, hash map[string]interface{}) (bool, error) {
	res, err := r.client.HMSet(ctx, key, hash).Result()
	return res, err
}

func (r *redis) HSet(ctx context.Context, key string, field string, value interface{}) (overwrite bool, err error) {
	res, err := r.client.HSet(ctx, key, field, value).Result()
	return res == 0, err
}

func (r *redis) HSetNX(ctx context.Context, key string, field string, value interface{}) (overwrite bool, err error) {
	return r.client.HSetNX(ctx, key, field, value).Result()
}

func (r *redis) HDel(ctx context.Context, key string, fields ...string) (delNum int64, err error) {
	return r.client.HDel(ctx, key, fields...).Result()
}

func (r *redis) HExists(ctx context.Context, key string, field string) (exists bool, err error) {
	return r.client.HExists(ctx, key, field).Result()
}

func (r *redis) HKeys(ctx context.Context, key string) ([]string, error) {
	return r.client.HKeys(ctx, key).Result()
}

func (r *redis) HLen(ctx context.Context, key string) (int64, error) {
	return r.client.HLen(ctx, key).Result()
}

func (r *redis) HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, field, incr).Result()
}

func (r *redis) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error) {
	return r.client.HIncrByFloat(ctx, key, field, incr).Result()
}
