package dispatch

import (
	"github.com/changsongl/delay-queue/bucket"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/pool"
	"github.com/changsongl/delay-queue/queue"
	"github.com/changsongl/delay-queue/type/job"
)

type Dispatch interface {
	Add(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (err error)
	Pop(topic job.Topic) (id job.Id, body job.Body, err error)
	Finish(topic job.Topic, id job.Id) (err error)
	Delete(topic job.Topic, id job.Id) (err error)
}

type dispatch struct {
	logger log.Logger
	bucket bucket.Bucket
	pool   pool.Pool
	queue  queue.Queue
}

func NewDispatch(logger log.Logger, new func() (bucket.Bucket, pool.Pool, queue.Queue)) Dispatch {
	b, p, q := new()

	return &dispatch{
		logger: logger,
		bucket: b,
		pool:   p,
		queue:  q,
	}
}

func (d dispatch) Add(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (err error) {
	panic("implement me")
}

func (d dispatch) Pop(topic job.Topic) (id job.Id, body job.Body, err error) {
	panic("implement me")
}

func (d dispatch) Finish(topic job.Topic, id job.Id) (err error) {
	panic("implement me")
}

func (d dispatch) Delete(topic job.Topic, id job.Id) (err error) {
	panic("implement me")
}
