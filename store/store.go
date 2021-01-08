package store

import "github.com/changsongl/delay-queue/job"

type Store interface {
	LoadJob() *job.Job
}
