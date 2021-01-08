package pool

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
)

type Pool interface {
	SaveJob(j *job.Job) error
}

type pool struct {
	s store.Store
}

func New(s store.Store) Pool {
	return pool{s: s}
}

// TODO: Save Job
func (p pool) SaveJob(j *job.Job) error {
	return nil
}
