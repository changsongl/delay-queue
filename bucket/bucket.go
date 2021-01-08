package bucket

import "github.com/changsongl/delay-queue/store"

type Bucket interface {
}

type bucket struct {
	s store.Store
}

func New(s store.Store) Bucket {
	return bucket{s: s}
}
