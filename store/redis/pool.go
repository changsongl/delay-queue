package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
)

func (s *storage) LoadJob(j *job.Job) error {
	jobData, err := s.rds.Get(context.Background(), j.GetName())
	if redis.IsError(err) {
		return err
	}

	return json.Unmarshal([]byte(jobData), j)
}

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
