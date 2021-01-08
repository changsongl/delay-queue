package bucket

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
)

type Bucket interface {
	SaveJob(j *job.Job) error
}

type bucket struct {
	s store.Store
}

func New(s store.Store) Bucket {
	return bucket{s: s}
}

// TODO: SaveJob
func (b bucket) SaveJob(j *job.Job) error {
	panic("implement me")
}
