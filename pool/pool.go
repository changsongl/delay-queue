package pool

import "github.com/changsongl/delay-queue/store"

type Pool interface {
}

type pool struct {
	s store.Store
}

func New(s store.Store) Pool {
	return pool{s: s}
}
