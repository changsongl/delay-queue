package redis

import (
	"context"
)

func (r *redis) SAdd(ctx context.Context, key string, member ...interface{}) (int64, error) {
	return r.Client.SAdd(ctx, key, member...).Result()
}

func (r *redis) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

func (r *redis) SRem(ctx context.Context, key string, member ...interface{}) (int64, error) {
	return r.Client.SRem(ctx, key, member...).Result()
}

func (r *redis) SPop(ctx context.Context, key string) (string, error) {
	return r.Client.SPop(ctx, key).Result()
}

func (r *redis) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	return r.Client.SRandMemberN(ctx, key, count).Result()
}

func (r *redis) SRandMember(ctx context.Context, key string) (string, error) {
	return r.Client.SRandMember(ctx, key).Result()
}

func (r *redis) SCard(ctx context.Context, key string) (int64, error) {
	return r.Client.SCard(ctx, key).Result()
}

func (r *redis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.Client.SIsMember(ctx, key, member).Result()
}
