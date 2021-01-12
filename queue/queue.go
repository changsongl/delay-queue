package queue

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
)

type Queue interface {
	Push(*job.Job) error
}

type queue struct {
	s    store.Store
	name string
}

func New(s store.Store, name string) Queue {
	return queue{s: s, name: name}
}

func (r queue) Push(j *job.Job) error {
	return r.s.PushJobToQueue(r.name, j)
}
