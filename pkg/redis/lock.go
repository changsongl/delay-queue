package redis

import (
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/go-redsync/redsync/v4"
	"time"
)

func (r *redis) GetLocker(name string) lock.Locker {
	return r.sync.NewMutex(name, redsync.WithExpiry(4*time.Second))
}
