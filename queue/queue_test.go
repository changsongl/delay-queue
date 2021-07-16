package queue

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	storemock "github.com/changsongl/delay-queue/test/mock/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueuePush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queueName := "test_queue_name"
	jobTopic := job.Topic("job_topic")
	queue := fmt.Sprintf("%s_%s", queueName, jobTopic)

	j, err := job.New(jobTopic, "1", 1, 1, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	sm := storemock.NewMockStore(ctrl)
	sm.EXPECT().PushJobToQueue(queue, j).Return(nil)
	q := New(sm, queueName)

	err = q.Push(j)
	require.NoError(t, err)
}

func TestQueuePop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queueName := "test_queue_name"
	jobTopic := job.Topic("job_topic")
	queue := fmt.Sprintf("%s_%s", queueName, jobTopic)

	expectNV := job.NameVersion("haha")
	sm := storemock.NewMockStore(ctrl)
	sm.EXPECT().PopJobFromQueue(queue).Return(expectNV, nil)
	q := New(sm, queueName)

	nv, err := q.Pop(jobTopic)
	require.NoError(t, err)
	require.Equal(t, expectNV, nv)
}
