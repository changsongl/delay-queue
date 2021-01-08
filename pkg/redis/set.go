package redis

import (
	"context"
)

func (r *redis) SAdd(ctx context.Context, key string, member ...interface{}) (int64, error) {
	return r.client.SAdd(ctx, key, member...).Result()
}

func (r *redis) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *redis) SRem(ctx context.Context, key string, member ...interface{}) (int64, error) {
	return r.client.SRem(ctx, key, member...).Result()
}

func (r *redis) SPop(ctx context.Context, key string) (string, error) {
	return r.client.SPop(ctx, key).Result()
}

func (r *redis) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	return r.client.SRandMemberN(ctx, key, count).Result()
}

func (r *redis) SRandMember(ctx context.Context, key string) (string, error) {
	return r.client.SRandMember(ctx, key).Result()
}

func (r *redis) SCard(ctx context.Context, key string) (int64, error) {
	return r.client.SCard(ctx, key).Result()
}

func (r *redis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}
