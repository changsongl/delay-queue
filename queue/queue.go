package queue

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
)

type Queue interface {
	AddJob(job.Job) bool
}

type queue struct {
	s store.Store
}

func New(s store.Store) Queue {
	return queue{s: s}
}

func (r queue) AddJob(j job.Job) bool {
	panic("implement me")
}
