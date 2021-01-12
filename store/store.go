package store

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
)

type Store interface {
	GetLock(name string) lock.Locker

	CreateJob(j *job.Job) error
	LoadJob(j *job.Job) error
	//DeleteJob(j *job.Job) error

	CreateJobInBucket(bucket string, j *job.Job) error
	GetReadyJobsInBucket(bucket string, num uint) ([]job.NameVersion, error)

	PushJobToQueue(queue string, j *job.Job) error
}
