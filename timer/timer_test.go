package timer

import (
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	logger, err := log.New()
	require.NoError(t, err)
	delay, interval := time.Second, 2*time.Second
	var noDelayCount, delayCount int64

	testCases := []struct {
		sleep   time.Duration
		atLeast int64
		task    TaskFunc
		count   *int64
	}{
		{sleep: 10 * time.Second, atLeast: 4, count: &noDelayCount, task: func() (hasMore bool, err error) {
			atomic.AddInt64(&noDelayCount, 1)
			return false, nil
		}},
		{sleep: 10 * time.Second, atLeast: 9, count: &delayCount, task: func() (hasMore bool, err error) {
			atomic.AddInt64(&delayCount, 1)
			return true, nil
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
