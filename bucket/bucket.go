package bucket

import (
	"errors"
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/store"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultMaxFetchNum uint64 = 20
)

var (
	metricOnce sync.Once
)

// Bucket interface to save jobs and repeat is searched
// for jobs which are ready to process
type Bucket interface {
	CreateJob(j *job.Job, isTTR bool) error
	GetBuckets() []uint64
	GetBucketJobs(bid uint64) ([]job.NameVersion, error)

	GetMaxFetchNum() uint64
	SetMaxFetchNum(num uint64)
}

// bucket implement Bucket interface
type bucket struct {
	s                store.Store          // real storage
	size             uint64               // bucket size for round robin
	name             string               // bucket name prefix
	count            *uint64              // current bucket
	locks            []lock.Locker        // locks for buckets
	maxFetchNum      uint64               // max number for fetching jobs
	l                log.Logger           // logger
	onFlightJobGauge *prometheus.GaugeVec // on flight jobs number in bucket
}

// New a Bucket interface object
func New(s store.Store, l log.Logger, size uint64, name string) Bucket {
	var c uint64 = 0
	b := &bucket{
		s:           s,
		size:        size,
		name:        name,
		count:       &c,
		l:           l,
		maxFetchNum: DefaultMaxFetchNum,
	}
	b.locks = make([]lock.Locker, 0, size)

	var i uint64 = 0
	for i < size {
		b.locks = append(b.locks, s.GetLock(b.getBucketNameById(i)))
		i++
	}

	b.CollectMetrics()

	return b
}

// CreateJob create job on bucket, bucket is selected
// by round robin policy
func (b *bucket) CreateJob(j *job.Job, isTTR bool) error {
	currentBucket := b.getNextBucket()
	err := b.s.CreateJobInBucket(currentBucket, j, isTTR)
	return err
}

// getNextBucket get next round robin bucket
func (b *bucket) getNextBucket() string {
	current := atomic.AddUint64(b.count, 1)
	return b.getBucketNameById(current % b.size)
}

// getBucketNameById return bucket name by id
func (b *bucket) getBucketNameById(id uint64) string {
	return fmt.Sprintf("%s_%d", b.name, id)
}

// GetBuckets return all bucket ids
func (b *bucket) GetBuckets() []uint64 {
	buckets := make([]uint64, 0, b.size)
	var i uint64 = 0
	for i < b.size {
		buckets = append(buckets, i)
		i++
	}

	return buckets
}

// GetBucketJobs return job.NameVersion which are ready to process. If this function
// call return names and the size of name is equal to num. Then it mean it may be
// more jobs are ready, but they are still in the bucket.
func (b *bucket) GetBucketJobs(bid uint64) ([]job.NameVersion, error) {
	bucketName := b.getBucketNameById(bid)
	nameVersions, err := b.s.GetReadyJobsInBucket(bucketName, uint(b.maxFetchNum))
	if err != nil {
		return nil, err
	}

	return nameVersions, nil
}

// GetLock return a lock for the given bucket id. use it when get jobs from bucket,
func (b *bucket) GetLock(bid uint64) (lock.Locker, error) {
	if bid >= b.size {
		return nil, errors.New("invalid bucket id to get lock")
	}

	return b.locks[bid], nil
}

// GetMaxFetchNum return the max number of job to fetch each time
func (b *bucket) GetMaxFetchNum() uint64 {
	return b.maxFetchNum
}

// SetMaxFetchNum set the max number of job to fetch each time
func (b *bucket) SetMaxFetchNum(num uint64) {
	b.maxFetchNum = num
}

func (b *bucket) CollectMetrics() {
	b.onFlightJobGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "delay_queue_in_flight_jobs_numbers_in_bucket",
		Help: "Gauge of the number of inflight jobs in each bucket",
	}, []string{"bucket"})

	metricOnce.Do(func() {
		err := prometheus.Register(b.onFlightJobGauge)
		if err != nil {
			b.l.Error("prometheus.Register failed", log.Error(err))
			return
		}
	})

	go func() {
		// TODO: graceful shutdown
		for {
			var i uint64
			for ; i < b.size; i++ {
				// collect
				bName := b.getBucketNameById(i)
				num, err := b.s.CollectInFlightJobNumberBucket(bName)
				if err != nil {
					b.l.Error("b.s.CollectInFlightJobNumberBucket failed", log.Error(err))
				}
				b.onFlightJobGauge.WithLabelValues(bName).Set(float64(num))
			}

			time.Sleep(30 * time.Second)
		}
	}()
}
