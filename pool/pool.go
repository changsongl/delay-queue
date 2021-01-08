package pool

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/encode"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/store"
)

type Pool interface {
	CreateJob(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) (*job.Job, error)
}

type pool struct {
	s store.Store
	l log.Logger
	e encode.Encoder
}

func New(s store.Store, l log.Logger) Pool {
	return pool{s: s, l: l}
}

// CreateJob lock the job and save job into storage
func (p pool) CreateJob(topic job.Topic, id job.Id,
	delay job.Delay, ttr job.TTR, body job.Body) (*job.Job, error) {

	j, err := job.New(topic, id, delay, ttr, body, p.s)
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
