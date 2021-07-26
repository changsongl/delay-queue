package redis

import (
	"context"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
	"time"
)

// PushJobToQueue push the job to the given redis queue
func (s *storage) PushJobToQueue(queuePrefix, queueName string, j *job.Job) error {
	ctx := context.Background()
	_, _ = s.rds.SAdd(ctx, queuePrefix, queueName)
	_, err := s.rds.LPush(ctx, queueName, j.GetNameWithVersion())
	return err
}

// PopJobFromQueue pop the job from redis queue
func (s *storage) PopJobFromQueue(queue string) (job.NameVersion, error) {
	nv, err := s.rds.RPop(context.Background(), queue)
	if redis.IsError(err) {
		return "", err
	} else if redis.IsNil(err) {
		return "", nil
	}
	return job.NameVersion(nv), nil
}

// BPopJobFromQueue pop the job from redis queue with block time
func (s *storage) BPopJobFromQueue(queue string, blockTime time.Duration) (job.NameVersion, error) {
	queueElement, err := s.rds.BRPop(context.Background(), queue, blockTime)
	if redis.IsError(err) {
		return "", err
	} else if redis.IsNil(err) {
		return "", nil
	} else if len(queueElement) != 2 {
		return "", nil
	}

	return job.NameVersion(queueElement[1]), nil
}

// CollectInFlightJobNumberQueue collect in flight job numbers in all queues
func (s *storage) CollectInFlightJobNumberQueue(queuePrefix string) (map[string]uint64, error) {
	ctx := context.Background()
	members, err := s.rds.SMembers(ctx, queuePrefix)
	if err != nil {
		return nil, err
	}

	queueMapNum := make(map[string]uint64, len(members))
	for _, member := range members {
		num, err := s.rds.LLen(ctx, member)
		if err != nil {
			return nil, err
		}
		queueMapNum[member] = uint64(num)
	}

	return queueMapNum, nil
}
