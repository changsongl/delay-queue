package queue

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
)

// Queue is a queue for ready jobs.
type Queue interface {
	Push(*job.Job) error
	Pop(topic job.Topic) (job.NameVersion, error)
}

// queue is Queue implementation struct.
type queue struct {
	s    store.Store
	name string
}

// New a queue with a name, and storage for queue.
func New(s store.Store, name string) Queue {
	return queue{s: s, name: name}
}

// Push a job to queue by job's topic
func (r queue) Push(j *job.Job) error {
	queue := r.getQueueName(j.Topic)
	return r.s.PushJobToQueue(queue, j)
}

// Pop a job from queue by job's topic
func (r queue) Pop(topic job.Topic) (job.NameVersion, error) {
	queue := r.getQueueName(topic)
	return r.s.PopJobFromQueue(queue)
}

// getQueueName based on topic of job and the queue name.
func (r queue) getQueueName(topic job.Topic) string {
	return fmt.Sprintf("%s_%s", r.name, topic)
}
