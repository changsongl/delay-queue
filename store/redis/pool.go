package redis

import (
	"context"
	"errors"
	"github.com/changsongl/delay-queue/job"
)

func (s *storage) LoadJob() (*job.Job, error) {
	return nil, nil
}

func (s *storage) CreateJob(j *job.Job) error {
	str, err := s.e.Encode(j)
	if err != nil {
		return err
	}

	result, err := s.r.SetNx(context.Background(), j.GetName(), str)
	if err != nil {
		return err
	} else if !result {
		return errors.New("job is exists")
	}

	return nil
}
