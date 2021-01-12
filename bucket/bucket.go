package bucket

import (
	"errors"
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/changsongl/delay-queue/store"
	"sync/atomic"
)

// Bucket interface to save jobs and repeat is searched
// for jobs which are ready to process
type Bucket interface {
	CreateJob(j *job.Job, isTTR bool) error
	GetBuckets() []uint64
	GetBucketJobs(bid uint64, num uint) ([]job.NameVersion, error)
}

// bucket implement Bucket interface
type bucket struct {
	s     store.Store   // real storage
	size  uint64        // bucket size for round robin
	name  string        // bucket name prefix
	count *uint64       // current bucket
	locks []lock.Locker // locks for buckets
}

// New a Bucket interface object
func New(s store.Store, size uint64, name string) Bucket {
	var c uint64 = 0
	b := &bucket{s: s, size: size, name: name, count: &c}
	b.locks = make([]lock.Locker, 0, size)

	var i uint64 = 0
	for i < size {
		b.locks = append(b.locks, s.GetLock(b.getBucketNameById(i)))
		i++
	}

	return b
}

// CreateJob create job on bucket, bucket is selected
// by round robin policy
func (b *bucket) CreateJob(j *job.Job, isTTR bool) error {
	currentBucket := b.getCurrentBucket()
	err := b.s.CreateJobInBucket(currentBucket, j, isTTR)
	return err
}

// getCurrentBucket get current round robin bucket
func (b *bucket) getCurrentBucket() string {
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
func (b *bucket) GetBucketJobs(bid uint64, num uint) ([]job.NameVersion, error) {
	bucketName := b.getBucketNameById(bid)
	nameVersions, err := b.s.GetReadyJobsInBucket(bucketName, num)
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
