package redis

import (
	"context"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
)

// PushJobToQueue push the job to the given redis queue
func (s *storage) PushJobToQueue(queue string, j *job.Job) error {
	_, err := s.rds.LPush(context.Background(), queue, j.GetNameWithVersion())
	return err
}

func (s *storage) PopJobFromQueue(queue string) (job.NameVersion, error) {
	jStr, err := s.rds.RPop(context.Background(), queue)
	if redis.IsError(err) {
		return "", err
	} else if redis.IsNil(err) {
		return "", nil
	}

	return job.NameVersion(jStr), nil
}
