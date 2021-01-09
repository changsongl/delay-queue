package store

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
)

type Store interface {
	GetLock(name string) lock.Locker

	CreateJob(*job.Job) error
	LoadJob() (*job.Job, error)

	CreateJobInBucket(bucket string, j *job.Job) error
	GetReadyJobsInBucket(bucket string, num uint) ([]job.NameVersion, error)
}
