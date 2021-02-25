package job

import (
	"errors"
	"fmt"
	"github.com/changsongl/delay-queue/pkg/lock"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	lockFunc = func(name string) lock.Locker {
		return &locker{}
	}
)

type locker struct {
	m sync.Mutex
}

func (l *locker) Unlock() (bool, error) {
	return true, nil
}

func (l *locker) Lock() error {
	return nil
}

func TestNew(t *testing.T) {
	delay, ttr, body := Delay(10*time.Second), TTR(5*time.Second), Body("")
	emptyErrMsg := errors.New("topic or id is empty")

	testCases := []struct {
		topic Topic
		id    Id
		err   error
	}{
		{topic: "", id: "", err: emptyErrMsg},
		{topic: "t1222", id: "id1", err: nil},
	}

	for _, tc := range testCases {
		job, err := New(tc.topic, tc.id, delay, ttr, body, lockFunc)
		if tc.err == nil {
			require.NoError(t, err, "job should have no error: %s", err)
			require.Equal(t, tc.topic, job.Topic, "wrong topic %s (expect: %s)", job.Topic, tc.topic)
			require.Equal(t, tc.id, job.ID, "wrong id %s (expect: %s)", job.ID, tc.id)
			require.Equal(t, delay, job.Delay, "wrong delay %s (expect: %s)", job.Delay, delay)
			require.Equal(t, ttr, job.TTR, "wrong ttr %s (expect: %s)", job.TTR, ttr)
			require.Equal(t, body, job.Body, "wrong body %s (expect: %s)", job.Body, body)
		} else {
			require.Error(t, err, "job should have error")
			require.Equal(t, err, tc.err, "job should have error %s (expect: %s)", err, tc.err)
		}
	}
}

func TestGet(t *testing.T) {
	emptyErrMsg := errors.New("topic or id is empty")
	testCases := []struct {
		topic Topic
		id    Id
		err   error
	}{
		{topic: "", id: "", err: emptyErrMsg},
		{topic: "t1222", id: "id1", err: nil},
	}

	for _, tc := range testCases {
		job, err := Get(tc.topic, tc.id, lockFunc)
		if tc.err == nil {
			require.NoError(t, err, "job should have no error: %s", err)
			require.Equal(t, tc.topic, job.Topic, "wrong topic %s (expect: %s)", job.Topic, tc.topic)
			require.Equal(t, tc.id, job.ID, "wrong id %s (expect: %s)", job.ID, tc.id)
		} else {
			require.Error(t, err, "job should have error")
			require.Equal(t, err, tc.err, "job should have error %s (expect: %s)", err, tc.err)
		}
	}
}

func TestIsVersionSame(t *testing.T) {
	job, err := New("topic", "id111", Delay(time.Second), TTR(time.Second), "", lockFunc)
	require.NoError(t, err, "Get should no error")

	same := job.IsVersionSame(job.Version)
	require.Equal(t, true, same, "version should be same")

	notSame := job.IsVersionSame(Version{t: time.Now()})
	require.Equal(t, false, notSame, "version should be not same")
}

func TestGetter(t *testing.T) {
	job, err := New("topic", "id111", Delay(time.Second), TTR(time.Second), "", lockFunc)
	require.NoError(t, err, "Get should no error")

	jName := fmt.Sprintf("%s_%s", job.Topic, job.ID)
	require.Equal(t, jName, job.GetName(), "job name should be same")
	require.EqualValues(t, fmt.Sprintf("%s_%s_%s", job.Topic, job.ID, job.Version), job.GetNameWithVersion(),
		"name with version should be same")
	require.EqualValues(t, fmt.Sprintf("%s_%s", jName, "lock"), job.getLockName(),
		"lock name should be same")

	// TODO: GetDelayTimeFromNow,GetTTRTimeFromNow need mock time.Now method
	err = job.Lock()
	require.NoError(t, err, "Lock should no error")

	result, err := job.Unlock()
	require.NoError(t, err, "Unlock should no error")
	require.EqualValues(t, true, result, "Unlock result should be true")
}
