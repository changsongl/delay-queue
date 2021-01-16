package store

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
)

// Store is a common storage interface, which is for manage jobs pool,
// buckets, and queue. Right now it is based on redis implementation,
// it might be change to other storage like memory or mysql, etc...
type Store interface {
	GetLock(name string) lock.Locker

	CreateJob(j *job.Job) error
	ReplaceJob(j *job.Job) error
	LoadJob(j *job.Job) error
	DeleteJob(j *job.Job) (bool, error)

	CreateJobInBucket(bucket string, j *job.Job, isTTR bool) error
	GetReadyJobsInBucket(bucket string, num uint) ([]job.NameVersion, error)

	PushJobToQueue(queue string, j *job.Job) error
	PopJobFromQueue(queue string) (job.NameVersion, error)
}
