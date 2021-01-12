package pool

import (
	"errors"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/encode"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/store"
)

type Pool interface {
	CreateJob(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (*job.Job, error)
	LoadReadyJob(topic job.Topic, id job.Id, version job.Version) (*job.Job, error)
}

type pool struct {
	s store.Store
	l log.Logger
	e encode.Encoder
}

func New(s store.Store, l log.Logger) Pool {
	return pool{s: s, l: l.WithModule("pool")}
}

// CreateJob lock the job and save job into storage
func (p pool) CreateJob(topic job.Topic, id job.Id,
	delay job.Delay, ttr job.TTR, body job.Body) (*job.Job, error) {

	j, err := job.New(topic, id, delay, ttr, body, p.s.GetLock)
	if err != nil {
		return nil, err
	}

	err = j.Lock()
	if err != nil {
		return nil, err
	}

	defer func() {
		if ok, err := j.Unlock(); !ok || err != nil {
			p.l.Error(
				"unlock failed",
				log.String("job", j.GetName()),
				log.Reflect("err", err),
			)
		}
	}()

	err = p.s.CreateJob(j)
	if err != nil {
		return nil, err
	}

	return j, err
}

// LoadReadyJob load ready job which is just gotten from bucket. this method will check
// job version is still same. If not same, then it means the just has been replaced, so
// this job should not process anymore.
func (p pool) LoadReadyJob(topic job.Topic, id job.Id, version job.Version) (*job.Job, error) {
	j, err := job.Get(topic, id, p.s.GetLock)
	if err != nil {
		return nil, err
	}

	err = p.s.LoadJob(j)
	if err != nil {
		return nil, err
	}

	if !j.IsVersionSame(version) {
		return nil, errors.New("version is not same")
	}

	return j, nil
}
