package redis

import (
	"github.com/changsongl/delay-queue/pkg/encode"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/changsongl/delay-queue/pkg/redis"
	"github.com/changsongl/delay-queue/store"
)

type storage struct {
	r redis.Redis
	e encode.Encoder
}

func (s storage) GetLock(name string) lock.Locker {
	return s.r.GetLocker(name)
}

func NewStore(r redis.Redis) store.Store {
	return &storage{r: r, e: encode.New()}
}
