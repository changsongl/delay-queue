package queue

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
)

type Queue interface {
	Push(*job.Job) error
	Pop(topic job.Topic) (job.NameVersion, error)
}

type queue struct {
	s    store.Store
	name string
}

func New(s store.Store, name string) Queue {
	return queue{s: s, name: name}
}

func (r queue) Push(j *job.Job) error {
	queue := r.getQueueName(j.Topic)
	return r.s.PushJobToQueue(queue, j)
}

func (r queue) Pop(topic job.Topic) (job.NameVersion, error) {
	queue := r.getQueueName(topic)
	return r.s.PopJobFromQueue(queue)
}

func (r queue) getQueueName(topic job.Topic) string {
	return fmt.Sprintf("%s_%s", r.name, topic)
}
