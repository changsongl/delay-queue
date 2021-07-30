package redis

import (
	"github.com/changsongl/delay-queue/pkg/encode"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/changsongl/delay-queue/pkg/redis"
	"github.com/changsongl/delay-queue/store"
)

// storage is store.Store implementation struct,
// it is using redis and encoder to save all data.
type storage struct {
	rds     redis.Redis
	encoder encode.Encoder
}

// GetLock based on given name, return a common locker.
func (s storage) GetLock(name string) lock.Locker {
	return s.rds.GetLocker(name)
}

// NewStore return a redis storage
func NewStore(r redis.Redis) store.Store {
	return &storage{rds: r, encoder: encode.NewCompress()}
}
