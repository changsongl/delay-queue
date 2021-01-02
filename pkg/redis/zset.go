package redis

import (
	"context"
	gredis "github.com/go-redis/redis/v8"
	"strconv"
)

type Z gredis.Z

func (r *redis) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.ZRevRange(ctx, key, start, stop).Result()
}

func (r *redis) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	zSli, err := r.Client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return []Z{}, err
	}

	results := make([]Z, 0, len(zSli))
	for _, z := range zSli {
		results = append(results, Z(z))
	}

	return results, nil
}

func (r *redis) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.ZRange(ctx, key, start, stop).Result()
}

func (r *redis) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	zSli, err := r.Client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return []Z{}, err
	}

	results := make([]Z, 0, len(zSli))
	for _, z := range zSli {
		results = append(results, Z(z))
	}

	return results, nil
}

func (r *redis) ZRangeByScoreWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	zSli, err := r.Client.ZRangeByScoreWithScores(ctx, key, &gredis.ZRangeBy{Min: startStr, Max: stopStr}).Result()
	if err != nil {
		return []Z{}, err
	}

	results := make([]Z, 0, len(zSli))
	for _, z := range zSli {
		results = append(results, Z(z))
	}

	return results, nil
}

func (r *redis) ZRangeByScoreWithScoresByOffset(ctx context.Context,
	key string, start, stop, offset, count int64) ([]Z, error) {
	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	zSli, err := r.Client.ZRangeByScoreWithScores(ctx, key,
		&gredis.ZRangeBy{Min: startStr, Max: stopStr, Offset: offset, Count: count}).Result()
	if err != nil {
		return []Z{}, err
	}

	results := make([]Z, 0, len(zSli))
	for _, z := range zSli {
		results = append(results, Z(z))
	}

	return results, nil
}

func (r *redis) ZRangeByScore(ctx context.Context, key string, start, stop int64) ([]string, error) {
	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	return r.Client.ZRangeByScore(ctx, key, &gredis.ZRangeBy{Min: startStr, Max: stopStr}).Result()
}

func (r *redis) ZRangeByScoreByOffset(ctx context.Context,
	key string, start, stop, offset, count int64) ([]string, error) {

	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	return r.Client.ZRangeByScore(ctx,
		key, &gredis.ZRangeBy{Min: startStr, Max: stopStr, Offset: offset, Count: count}).Result()
}

func (r *redis) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return r.Client.ZRevRank(ctx, key, member).Result()
}

func (r *redis) ZRevRangeByScore(ctx context.Context, key string, start, stop int64) ([]string, error) {
	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	return r.Client.ZRevRangeByScore(ctx, key, &gredis.ZRangeBy{Min: startStr, Max: stopStr}).Result()
}

func (r *redis) ZRevRangeByScoreByOffset(ctx context.Context,
	key string, start, stop, offset, count int64) ([]string, error) {

	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	return r.Client.ZRevRangeByScore(ctx,
		key, &gredis.ZRangeBy{Min: startStr, Max: stopStr, Offset: offset, Count: count}).Result()
}

func (r *redis) ZRevRangeByScoreWithScores(ctx context.Context,
	key string, start, stop int64) ([]gredis.Z, error) {

	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	res, err := r.Client.ZRevRangeByScoreWithScores(ctx,
		key, &gredis.ZRangeBy{Min: startStr, Max: stopStr}).Result()
	if err != nil && err != gredis.Nil {
		return []gredis.Z{}, err
	}

	return res, nil
}

func (r *redis) ZRevRangeByScoreWithScoresByOffset(ctx context.Context,
	key string, start, stop, offset, count int64) ([]gredis.Z, error) {

	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	res, err := r.Client.ZRevRangeByScoreWithScores(ctx,
		key, &gredis.ZRangeBy{Min: startStr, Max: stopStr, Offset: offset, Count: count}).Result()
	if err != nil && err != gredis.Nil {
		return []gredis.Z{}, err
	}

	return res, nil
}

func (r *redis) ZCard(ctx context.Context, key string) (int64, error) {
	return r.Client.ZCard(ctx, key).Result()
}

func (r *redis) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return r.Client.ZScore(ctx, key, member).Result()
}

func (r *redis) ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	if len(members) == 0 {
		return 0, nil
	}

	addSli := make([]*gredis.Z, 0, len(members))
	for _, member := range members {
		m := gredis.Z(member)
		addSli = append(addSli, &m)
	}

	return r.Client.ZAdd(ctx, key, addSli...).Result()
}

func (r *redis) ZCount(ctx context.Context, key string, start, stop int64) (int64, error) {
	startStr, stopStr := strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)
	return r.Client.ZCount(ctx, key, startStr, stopStr).Result()
}

func (r *redis) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.Client.ZRem(ctx, key, members...).Result()
}

func (r *redis) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return r.Client.ZRemRangeByRank(ctx, key, start, stop).Result()
}
