package dispatch

import (
	"github.com/changsongl/delay-queue/bucket"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/pool"
	"github.com/changsongl/delay-queue/queue"
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
	// create job and save to pool with lock
	j, err := job.New(topic, id, delay, ttr, body)
	if err != nil {
		return err
	}

	err = j.Lock()
	if err != nil {
		return err
	}

	defer func() {
		errDefer := j.Unlock()
		if err == nil {
			err = errDefer
		}
	}()

	err = d.pool.SaveJob(j)
	if err != nil {
		return err
	}

	// choose bucket to save
	err = d.bucket.SaveJob(j)
	return err
}

func (d dispatch) Pop(topic job.Topic) (id job.Id, body job.Body, err error) {
	// find job from ready queue

	// if get a message, check it is valid

	// prepare for ttr time

	return "", "", nil
}

func (d dispatch) Finish(topic job.Topic, id job.Id) (err error) {
	// find job

	// set it is done

	return nil
}

func (d dispatch) Delete(topic job.Topic, id job.Id) (err error) {
	// find job

	// check job current status

	// delete if ok

	return nil
}
