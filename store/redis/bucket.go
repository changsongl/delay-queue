package redis

import (
	"context"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
)

// CreateJobInBucket create job which will be ready after delay time
func (s *storage) CreateJobInBucket(bucketName string, j *job.Job) error {
	delayTime := float64(j.GetDelayTimeFromNow().Unix())
	_, err := s.r.ZAdd(context.Background(), bucketName, redis.Z{
		Score:  delayTime,
		Member: j.GetName(),
	})

	return err
}
