package redis

import (
	"github.com/changsongl/delay-queue/pkg/encode"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/changsongl/delay-queue/pkg/redis"
	"github.com/changsongl/delay-queue/store"
)

type storage struct {
	rds     redis.Redis
	encoder encode.Encoder
}

func (s storage) GetLock(name string) lock.Locker {
	return s.rds.GetLocker(name)
}

func NewStore(r redis.Redis) store.Store {
	return &storage{rds: r, encoder: encode.New()}
}
