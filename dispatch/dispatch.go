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
	Add(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body, override bool) (err error)
	Pop(topic job.Topic) (j *job.Job, err error)
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
	d.timer.AddTask(func() (bool, error) {
		nameVersions, err := d.bucket.GetBucketJobs(bid)
		if err != nil {
			d.logger.Error("timer.task bucket.GetBucketJobs failed", log.String("err", err.Error()))
			return false, err
		}

		for _, nameVersion := range nameVersions {
			d.logger.Debug("process", log.String("nameVersion", string(nameVersion)))
			topic, id, version, err := nameVersion.Parse()
			if err != nil {
				d.logger.Error("timer.task nameVersion.Parse failed",
					log.String("err", err.Error()), log.String("nameVersion", string(nameVersion)))
				continue
			}

			j, err := d.pool.LoadReadyJob(topic, id, version)
			if err != nil {
				d.logger.Error("timer.task pool.LoadReadyJob failed",
					log.String("err", err.Error()), log.String("topic", string(topic)),
					log.String("id", string(id)), log.String("version", version.String()))
				continue
			}

			err = d.queue.Push(j)
			if err != nil {
				d.logger.Error("timer.task queue.Push failed", log.String("err", err.Error()))
			}
		}

		return len(nameVersions) == int(d.bucket.GetMaxFetchNum()), nil
	})
}

// Add job to job pool and push to bucket.
func (d dispatch) Add(topic job.Topic, id job.Id,
	delay job.Delay, ttr job.TTR, body job.Body, override bool) (err error) {

	j, err := d.pool.CreateJob(topic, id, delay, ttr, body, override)
	if err != nil {
		return err
	}

	return d.bucket.CreateJob(j, false)
}

// Pop job from bucket and return job info to let user process. if the ttr time
// is not zero, it will requeue after ttr time. if user doesn't call finish before
// that time, then this job can be pop again. User need to make sure ttr time is
// reasonable.
func (d dispatch) Pop(topic job.Topic) (j *job.Job, err error) {

	// find job from ready queue
	nameVersion, err := d.queue.Pop(topic)
	if err != nil {
		return
	} else if nameVersion == "" {
		err = nil
		return
	}

	topic, id, version, err := nameVersion.Parse()
	if err != nil {
		return
	}

	j, err = d.pool.LoadReadyJob(topic, id, version)
	if err != nil {
		return
	}

	if j.TTR != 0 {
		err := d.bucket.CreateJob(j, true)
		if err != nil {
			d.logger.Error("bucket ttr requeue failed", log.String("err", err.Error()))
		}
	}

	return j, nil
}

// Finish job. ack the processed job after user has done their job.
// delay queue will stop retrying and delete all information.
func (d dispatch) Finish(topic job.Topic, id job.Id) (err error) {
	// set it is done
	return d.pool.DeleteJob(topic, id)
}

// Delete job before. only delete job, when the bucket event is trigger,
// it gonna find the job is deleted, so it won't push to the ready queue.
func (d dispatch) Delete(topic job.Topic, id job.Id) (err error) {
	// delete job
	return d.pool.DeleteJob(topic, id)
}
