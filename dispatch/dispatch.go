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
		logger: logger.WithModule("dispatch"),
		bucket: b,
		pool:   p,
		queue:  q,
		timer:  t,
	}
}

// Run the dispatch with timer for getting ready jobs
func (d dispatch) Run() {
	buckets := d.bucket.GetBuckets()

	for _, b := range buckets {
		d.addTask(b)
	}

	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
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

// addTask the task is to get ready jobs from bucket and check data is valid,
// if yes then push to ready queue, if not then discard.
func (d dispatch) addTask(bid uint64) {
	d.timer.AddTask(func(num int) (int, error) {
		nameVersions, err := d.bucket.GetBucketJobs(bid, uint(num))
		if err != nil {
			d.logger.Error("task failed", log.String("err", err.Error()))
			return 0, err
		}

		for _, nameVersion := range nameVersions {
			d.logger.Debug("process", log.String("nameVersion", string(nameVersion)))
			topic, id, version, err := nameVersion.Parse()
			if err != nil {
				d.logger.Error("nameVersion.Parse failed", log.String("err", err.Error()))
				continue
			}

			j, err := d.pool.LoadReadyJob(topic, id, version)
			if err != nil {
				d.logger.Error("pool.LoadReadyJob failed", log.String("err", err.Error()))
				continue
			}

			err = d.queue.Push(j)
			if err != nil {
				d.logger.Error("queue.Push failed", log.String("err", err.Error()))
			}
		}

		return len(nameVersions), nil
	})
}

// Add job to job pool and push to bucket.
func (d dispatch) Add(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (err error) {
	j, err := d.pool.CreateJob(topic, id, delay, ttr, body)
	if err != nil {
		return err
	}

	return d.bucket.CreateJob(j)
}

// Pop job from bucket and return job info to let user process.
func (d dispatch) Pop(topic job.Topic) (id job.Id, body job.Body, err error) {
	// find job from ready queue

	// if get a message, check it is valid

	// prepare for ttr time

	return "", "", nil
}

// Finish job. ack the processed job after user has done their job.
// delay queue will stop retrying and delete all information.
func (d dispatch) Finish(topic job.Topic, id job.Id) (err error) {
	// find job

	// set it is done

	return nil
}

// Delete job before. only delete job, when the bucket event is trigger,
// it gonna find the job is deleted, so it won't push to the ready queue.
func (d dispatch) Delete(topic job.Topic, id job.Id) (err error) {
	// find job

	// check job current status

	// delete if ok

	return nil
}
