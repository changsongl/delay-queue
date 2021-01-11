package dispatch

import (
	"github.com/changsongl/delay-queue/bucket"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/pool"
	"github.com/changsongl/delay-queue/queue"
	"github.com/changsongl/delay-queue/timer"
	"os"
	"os/signal"
	"syscall"
)

type Dispatch interface {
	Add(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (err error)
	Pop(topic job.Topic) (id job.Id, body job.Body, err error)
	Finish(topic job.Topic, id job.Id) (err error)
	Delete(topic job.Topic, id job.Id) (err error)

	Run()
}

type dispatch struct {
	logger log.Logger
	bucket bucket.Bucket
	pool   pool.Pool
	queue  queue.Queue
	timer  timer.Timer
}

func NewDispatch(logger log.Logger, new func() (bucket.Bucket, pool.Pool, queue.Queue, timer.Timer)) Dispatch {
	b, p, q, t := new()

	return &dispatch{
		logger: logger,
		bucket: b,
		pool:   p,
		queue:  q,
		timer:  t,
	}
}

func (d dispatch) Run() {
	buckets := d.bucket.GetBuckets()

	for _, b := range buckets {
		d.addTask(b)
	}

	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-term:
				d.logger.Info("Signal dispatch stop")
				d.timer.Close()
				return
			}
		}
	}()

	d.logger.Info("Run dispatch")
	d.timer.Run()
	d.logger.Info("Dispatch is stopped")
}

func (d dispatch) addTask(bid uint64) {
	d.timer.AddTask(func(num int) (int, error) {
		nameVersions, err := d.bucket.GetBucketJobs(bid, uint(num))
		if err != nil {
			d.logger.Error("task failed", log.String("err", err.Error()))
			return 0, err
		}

		for _, nameVersion := range nameVersions {
			d.logger.Info("process", log.String("nameVersion", string(nameVersion)))
		}

		return len(nameVersions), nil
	})
}

func (d dispatch) Add(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (err error) {
	j, err := d.pool.CreateJob(topic, id, delay, ttr, body)
	if err != nil {
		return err
	}

	return d.bucket.CreateJob(j)
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
