package redis

import (
	"context"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/redis"
	"time"
)

// CreateJobInBucket create job which will be ready after delay time
func (s *storage) CreateJobInBucket(bucketName string, j *job.Job, isTTR bool) error {
	var delayTime float64
	if isTTR {
		delayTime = float64(j.GetTTRTimeFromNow().Unix())
	} else {
		delayTime = float64(j.GetDelayTimeFromNow().Unix())
	}

	_, err := s.rds.ZAdd(context.Background(), bucketName, redis.Z{
		Score:  delayTime,
		Member: j.GetNameWithVersion(),
	})

	return err
}

// GetReadyJobsInBucket get job which is ready to be pushed to queue
func (s *storage) GetReadyJobsInBucket(bucket string, num uint) ([]job.NameVersion, error) {
	nameStrings, err := s.rds.ZRangeByScoreByOffset(
		context.Background(),
		bucket,
		0,
		time.Now().Unix(),
		0,
		int64(num),
	)

	nvs := make([]job.NameVersion, 0, len(nameStrings))

	if err != nil {
		return nil, err
	} else if len(nameStrings) == 0 {
		return nvs, nil
	}

	_, err = s.rds.ZRem(context.Background(), bucket, redis.StringMembersToInterface(nameStrings)...)
	if err != nil {
		return nil, err
	}

	for _, nameString := range nameStrings {
		nvs = append(nvs, job.NewNameVersion(nameString))
	}
	return nvs, nil
}
