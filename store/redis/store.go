package redis

import (
	"github.com/changsongl/delay-queue/pkg/redis"
	"github.com/changsongl/delay-queue/store"
)

type storage struct {
	r redis.Redis
}

func NewStore(r redis.Redis) store.Store {
	return storage{r: r}
}
