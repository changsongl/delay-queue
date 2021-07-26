package queue

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	logmock "github.com/changsongl/delay-queue/test/mock/log"
	storemock "github.com/changsongl/delay-queue/test/mock/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueuePush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queueName := "test_queue_name"
	jobTopic := job.Topic("job_topic")
	que := fmt.Sprintf("%s_%s", queueName, jobTopic)

	j, err := job.New(jobTopic, "1", 1, 1, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	sm := storemock.NewMockStore(ctrl)
	mLog := logmock.NewMockLogger(ctrl)

	sm.EXPECT().PushJobToQueue(queueName, que, j).Return(nil)
	sm.EXPECT().CollectInFlightJobNumberQueue(queueName).AnyTimes()
	q := New(sm, mLog, queueName, time.Duration(0))

	err = q.Push(j)
	require.NoError(t, err)
}

func TestQueuePop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queueName := "test_queue_name"
	jobTopic := job.Topic("job_topic")
	que := fmt.Sprintf("%s_%s", queueName, jobTopic)

	expectNV := job.NameVersion("haha")
	sm := storemock.NewMockStore(ctrl)
	mLog := logmock.NewMockLogger(ctrl)
	blockTime := time.Duration(0)

	sm.EXPECT().PopJobFromQueue(que, blockTime).Return(expectNV, nil)
	sm.EXPECT().CollectInFlightJobNumberQueue(queueName).AnyTimes()
	q := New(sm, mLog, queueName, blockTime)

	nv, err := q.Pop(jobTopic)
	require.NoError(t, err)
	require.Equal(t, expectNV, nv)
}

func TestQueuePopWithBlockTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queueName := "test_queue_name"
	jobTopic := job.Topic("job_topic")
	que := fmt.Sprintf("%s_%s", queueName, jobTopic)

	expectNV := job.NameVersion("haha")
	sm := storemock.NewMockStore(ctrl)
	mLog := logmock.NewMockLogger(ctrl)
	blockTime := 2 * time.Second

	sm.EXPECT().PopJobFromQueue(que, blockTime).DoAndReturn(
		func(queue string, blockTime time.Duration) (job.NameVersion, error) {
			time.Sleep(blockTime)
			return expectNV, nil
		},
	)
	sm.EXPECT().CollectInFlightJobNumberQueue(queueName).AnyTimes()
	q := New(sm, mLog, queueName, blockTime)

	startTime := time.Now()
	nv, err := q.Pop(jobTopic)
	dur := time.Since(startTime)

	require.NoError(t, err)
	require.Equal(t, expectNV, nv)
	require.GreaterOrEqual(t, dur, blockTime)
}
