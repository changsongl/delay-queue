package redis

import (
	"context"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
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
	jStr, err := s.rds.RPop(context.Background(), queue)
	if redis.IsError(err) {
		return "", err
	} else if redis.IsNil(err) {
		return "", nil
	}

	return job.NameVersion(jStr), nil
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
