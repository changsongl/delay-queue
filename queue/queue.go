package queue

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/store"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"time"
)

// Queue is a queue for ready jobs.
type Queue interface {
	Push(*job.Job) error
	Pop(topic job.Topic) (job.NameVersion, error)
}

// queue is Queue implementation struct.
type queue struct {
	s                store.Store          // storage
	prefix           string               // prefix
	onFlightJobGauge *prometheus.GaugeVec // on flight jobs number in bucket
	l                log.Logger           // logger
}

// New a queue with a prefix, and storage for queue.
func New(s store.Store, l log.Logger, name string) Queue {
	q := &queue{s: s, prefix: name, l: l}
	q.CollectMetrics()
	return q
}

// Push a job to queue by job's topic
func (r *queue) Push(j *job.Job) error {
	que := r.getQueueName(j.Topic)
	prefix := r.getQueuePrefix()
	return r.s.PushJobToQueue(prefix, que, j)
}

// Pop a job from queue by job's topic
func (r *queue) Pop(topic job.Topic) (job.NameVersion, error) {
	que := r.getQueueName(topic)
	return r.s.PopJobFromQueue(que)
}

// getQueueName based on topic of job and the queue prefix.
func (r *queue) getQueueName(topic job.Topic) string {
	return fmt.Sprintf("%s_%s", r.prefix, topic)
}

// getQueuePrefix get queue prefix
func (r *queue) getQueuePrefix() string {
	return r.prefix
}

// CollectMetrics collect prometheus metrics for in flight jobs in queue
func (r *queue) CollectMetrics() {
	r.onFlightJobGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "delay_queue_in_flight_jobs_numbers_in_queue",
		Help: "Gauge of the number of inflight jobs in each queue",
	}, []string{"topic"})

	err := prometheus.Register(r.onFlightJobGauge)
	if err != nil {
		r.l.Error("prometheus.Register failed", log.Error(err))
		return
	}

	go func() {
		// TODO: graceful shutdown
		for {
			queueMapJobNum, err := r.s.CollectInFlightJobNumberQueue(r.prefix)
			if err != nil {
				r.l.Error("b.s.CollectInFlightJobNumberBucket failed", log.Error(err))
			}

			for queueName, num := range queueMapJobNum {
				topicName := strings.TrimLeft(queueName, r.getQueuePrefix())
				r.onFlightJobGauge.WithLabelValues(topicName).Set(float64(num))
			}

			time.Sleep(30 * time.Second)
		}
	}()
}
