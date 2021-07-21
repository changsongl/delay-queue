package redis

import (
	"context"
	"github.com/changsongl/delay-queue/config"
	"github.com/changsongl/delay-queue/pkg/lock"
	gredis "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"time"
)

// TODO: add pipeline

type Redis interface {
	Del(ctx context.Context, key string) (bool, error)
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)
	Ttl(ctx context.Context, key string) (time.Duration, error)
	FlushDB(ctx context.Context) error

	// kv
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	SetEx(ctx context.Context, key string, value interface{}, expire time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, increment int64) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, decrement int64) (int64, error)
	MGet(ctx context.Context, keys ...string) ([]*string, error)
	MSet(ctx context.Context, kvs map[string]interface{}) error
	SetNx(ctx context.Context, key string, value interface{}) (bool, error)
	SetNxExpire(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	// hash
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HMGet(ctx context.Context, key string, fields []string) ([]*string, error)
	HGet(ctx context.Context, key string, field string) (string, error)
	HMSet(ctx context.Context, key string, hash map[string]interface{}) (bool, error)
	HSet(ctx context.Context, key string, field string, value interface{}) (overwrite bool, err error)
	HSetNX(ctx context.Context, key string, field string, value interface{}) (overwrite bool, err error)
	HDel(ctx context.Context, key string, fields ...string) (delNum int64, err error)
	HExists(ctx context.Context, key string, field string) (exists bool, err error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HLen(ctx context.Context, key string) (int64, error)
	HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error)
	HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error)

	// zset
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error)
	ZRangeByScoreWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error)
	ZRangeByScoreWithScoresByOffset(ctx context.Context, key string, start, stop, offset, count int64) ([]Z, error)
	ZRangeByScore(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeByScoreByOffset(ctx context.Context, key string, start, stop, offset, count int64) ([]string, error)
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
	ZRevRangeByScore(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRangeByScoreByOffset(ctx context.Context, key string, start, stop, offset, count int64) ([]string, error)
	ZRevRangeByScoreWithScores(ctx context.Context, key string, start, stop int64) ([]gredis.Z, error)
	ZRevRangeByScoreWithScoresByOffset(ctx context.Context, key string, start, stop, offset, count int64) ([]gredis.Z, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	ZAdd(ctx context.Context, key string, members ...Z) (int64, error)
	ZCount(ctx context.Context, key string, start, stop int64) (int64, error)
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error)

	// list
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPop(ctx context.Context, key string) (string, error)
	LPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	LIndex(ctx context.Context, key string, idx int64) (string, error)
	LTrim(ctx context.Context, key string, start, stop int64) (string, error)

	// set
	SAdd(ctx context.Context, key string, member ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, member ...interface{}) (int64, error)
	SPop(ctx context.Context, key string) (string, error)
	SRandMemberN(ctx context.Context, key string, count int64) ([]string, error)
	SRandMember(ctx context.Context, key string) (string, error)
	SCard(ctx context.Context, key string) (int64, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	GetLocker(name string) lock.Locker

	Close() (err error)
}

type redis struct {
	client *gredis.Client
	sync   *redsync.Redsync
}

func New(conf config.Redis) Redis {
	cli := gredis.NewClient(
		&gredis.Options{
			Network:      conf.Network,
			Addr:         conf.Address,
			Username:     conf.Username,
			Password:     conf.Password,
			DB:           conf.DB,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
		},
	)

	rs := redsync.New(goredis.NewPool(cli))

	return &redis{
		client: cli,
		sync:   rs,
	}
}

// Close close client and all connections
func (r *redis) Close() (err error) {
	if r.client != nil {
		err = r.client.Close()
	}
	return
}
