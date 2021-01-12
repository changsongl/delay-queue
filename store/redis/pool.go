package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
)

// LoadJob information
func (s *storage) LoadJob(j *job.Job) error {
	jobData, err := s.rds.Get(context.Background(), j.GetName())
	if redis.IsError(err) {
		return err
	} else if redis.IsNil(err) {
		return errors.New("job is not exists")
	}

	return json.Unmarshal([]byte(jobData), j)
}

// CreateJob information, only if the job is not exists.
func (s *storage) CreateJob(j *job.Job) error {
	str, err := s.encoder.Encode(j)
	if err != nil {
		return err
	}

	result, err := s.rds.SetNx(context.Background(), j.GetName(), str)
	if err != nil {
		return err
	} else if !result {
		return errors.New("job is exists")
	}

	return nil
}

// DeleteJob delete job in redis
func (s *storage) DeleteJob(j *job.Job) (bool, error) {
	return s.rds.Del(context.Background(), j.GetName())
}
