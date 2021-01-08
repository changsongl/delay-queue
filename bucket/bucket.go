package bucket

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/store"
	"sync/atomic"
)

// Bucket interface to save jobs and repeat is searched
// for jobs which are ready to process
type Bucket interface {
	CreateJob(j *job.Job) error
}

// bucket implement Bucket interface
type bucket struct {
	s     store.Store // real storage
	size  uint64      // bucket size for round robin
	name  string      // bucket name prefix
	count uint64      // current bucket
}

// New a Bucket interface object
func New(s store.Store, size uint64, name string) Bucket {
	return bucket{s: s, size: size, name: name, count: 0}
}

// CreateJob create job on bucket, bucket is selected
// by round robin policy
func (b bucket) CreateJob(j *job.Job) error {
	currentBucket := b.GetCurrentBucket()
	err := b.s.CreateJobInBucket(currentBucket, j)
	return err
}

// GetCurrentBucket get current round robin bucket
func (b bucket) GetCurrentBucket() string {
	current := atomic.AddUint64(&b.count, 1)
	return fmt.Sprintf("%s_%d", b.name, current%b.size)
}
