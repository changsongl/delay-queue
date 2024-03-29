package pool

import (
	"errors"

	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/store"
)

var (
	// ErrJobNotExist error job is not exists
	ErrJobNotExist = errors.New("job is not exists")
	// ErrVersionNotSame error job version is not same
	ErrVersionNotSame = errors.New("version is not same")
)

// Pool is an interface for manage information of jobs.
type Pool interface {
	CreateJob(topic job.Topic, id job.ID,
		delay job.Delay, ttr job.TTR, body job.Body, override bool) (*job.Job, error)
	LoadReadyJob(topic job.Topic, id job.ID, version job.Version) (*job.Job, error)
	DeleteJob(topic job.Topic, id job.ID) error
}

// pool is Pool implementation struct
type pool struct {
	s store.Store
	l log.Logger
}

// New a pool with logger and storage
func New(s store.Store, l log.Logger) Pool {
	return pool{s: s, l: l.WithModule("pool")}
}

// CreateJob lock the job and save job into storage
func (p pool) CreateJob(topic job.Topic, id job.ID,
	delay job.Delay, ttr job.TTR, body job.Body, override bool) (*job.Job, error) {

	j, err := job.New(topic, id, delay, ttr, body, p.s.GetLock)
	if err != nil {
		return nil, err
	}

	err = j.Lock()
	if err != nil {
		return nil, err
	}

	defer func() {
		if ok, unlockErr := j.Unlock(); !ok || unlockErr != nil {
			p.l.Error(
				"j.Unlock failed",
				log.String("job", j.GetName()),
				log.Error(unlockErr),
				log.Bool("ok", ok),
			)
		}
	}()

	if override {
		err = p.s.ReplaceJob(j)
	} else {
		err = p.s.CreateJob(j)
	}

	if err != nil {
		return nil, err
	}

	return j, err
}

// LoadReadyJob load ready job which is just gotten from bucket. this method will check
// job version is still same. If not same, then it means the just has been replaced, so
// this job should not process anymore.
func (p pool) LoadReadyJob(topic job.Topic, id job.ID, version job.Version) (*job.Job, error) {
	j, err := job.Get(topic, id, p.s.GetLock)
	if err != nil {
		return nil, err
	}

	err = p.s.LoadJob(j)
	if err != nil {
		return nil, err
	}

	if !j.IsVersionSame(version) {
		return nil, ErrVersionNotSame
	}

	return j, nil
}

// DeleteJob a job, it will prevent job to be send to user.
func (p pool) DeleteJob(topic job.Topic, id job.ID) error {
	j, err := job.Get(topic, id, p.s.GetLock)
	if err != nil {
		return err
	}

	err = j.Lock()
	if err != nil {
		return err
	}

	defer func() {
		if ok, unlockErr := j.Unlock(); !ok || unlockErr != nil {
			p.l.Error(
				"j.Unlock failed",
				log.String("job", j.GetName()),
				log.Error(unlockErr),
				log.Bool("ok", ok),
			)
		}
	}()

	result, err := p.s.DeleteJob(j)
	if err != nil {
		return err
	} else if !result {
		return ErrJobNotExist
	}

	return nil
}
