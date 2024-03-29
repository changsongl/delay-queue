package integration

import (
	"context"
	"fmt"
	"github.com/changsongl/delay-queue/config"
	"github.com/changsongl/delay-queue/pkg/redis"
)

const (
	// RedisAddr redis default address
	RedisAddr = "127.0.0.1:6379"
	// DelayQueueAddr delay queue default address
	DelayQueueAddr = "http://127.0.0.1:8000"
)

var redisInstance redis.Redis

func init() {
	if err := CleanTestingStates(); err != nil {
		panic(fmt.Sprintf("Integration test failed: init(): %s", err.Error()))
	}
}

// CleanTestingStates clean all states from the previous testing
func CleanTestingStates() error {
	return GetRedis().FlushDB(context.Background())
}

// GetRedis get redis
func GetRedis() redis.Redis {
	if redisInstance == nil {
		redisConf := config.New().Redis
		redisConf.Address = RedisAddr
		redisInstance = redis.New(config.New().Redis)
	}

	return redisInstance
}

// AddJobRecord add job to set
func AddJobRecord(key, job string) error {
	_, err := GetRedis().SAdd(context.Background(), key, job)
	return err
}

// DeleteJobRecord delete job record
func DeleteJobRecord(key, job string) error {
	_, err := GetRedis().SRem(context.Background(), key, job)
	return err
}

// RecordNumbers get record num
func RecordNumbers(key string) (int64, error) {
	return GetRedis().SCard(context.Background(), key)
}
