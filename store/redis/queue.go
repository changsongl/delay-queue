package redis

import (
	"context"
	"github.com/changsongl/delay-queue/job"
)

// PushJobToQueue push the job to the given redis queue
func (s *storage) PushJobToQueue(queue string, j *job.Job) error {
	_, err := s.rds.RPush(context.Background(), queue, j.GetNameWithVersion())
	return err
}
