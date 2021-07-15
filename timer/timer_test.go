package timer

import (
	"errors"
	log "github.com/changsongl/delay-queue/test/mock/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.NewMockLogger(ctrl)
	logger.EXPECT().WithModule(gomock.Any()).MaxTimes(3).Return(logger)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).MinTimes(1)

	delay, interval := time.Second, 2*time.Second
	var noDelayCount, delayCount int64
	testCases := []struct {
		sleep   time.Duration
		atLeast int64
		task    TaskFunc
		count   *int64
	}{
		{sleep: 5 * interval, atLeast: 4, count: &noDelayCount, task: func() (hasMore bool, err error) {
			atomic.AddInt64(&noDelayCount, 1)
			return false, nil
		}},
		{sleep: 5 * interval, atLeast: 9, count: &delayCount, task: func() (hasMore bool, err error) {
			atomic.AddInt64(&delayCount, 1)
			return true, nil
		}},
		{sleep: 2 * interval, atLeast: 1, count: &delayCount, task: func() (hasMore bool, err error) {
			atomic.AddInt64(&delayCount, 1)
			return false, errors.New("error")
		}},
	}

	for _, test := range testCases {
		t.Log("start TestTimer case")
		tm := New(logger, interval, delay)
		tm.AddTask(test.task)
		t.Log("timer run")
		go tm.Run()

		time.Sleep(test.sleep)

		checkNum := atomic.LoadInt64(test.count)
		require.GreaterOrEqual(t, checkNum, test.atLeast)

		t.Log("timer close")
		tm.Close()
		t.Log("end TestTimer case")
	}
}
