package pool

import (
	"errors"
	"github.com/agiledragon/gomonkey"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	log "github.com/changsongl/delay-queue/test/mock/log"
	mocklock "github.com/changsongl/delay-queue/test/mock/pkg/lock"
	store "github.com/changsongl/delay-queue/test/mock/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateJobNewJobErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	mockLogger := log.NewMockLogger(ctrl)
	mockLogger.EXPECT().WithModule(gomock.Any()).Return(mockLogger)
	p := New(mockStore, mockLogger)

	expJob := &job.Job{}
	expErr := errors.New("job error")

	gomonkey.ApplyFunc(job.New, func(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR,
		body job.Body, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return expJob, expErr
	})

	j, err := p.CreateJob("", "", 1, 1, "", true)
	require.Equal(t, expErr, err)
	require.Nil(t, j)
}

func TestCreateJobLockErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)

	mockLogger := log.NewMockLogger(ctrl)
	mockLogger.EXPECT().WithModule(gomock.Any()).Return(mockLogger)

	mockLock := mocklock.NewMockLocker(ctrl)
	expJob := &job.Job{Mutex: mockLock}

	expErr := errors.New("lock error")
	mockLock.EXPECT().Lock().Return(expErr)

	p := New(mockStore, mockLogger)

	gomonkey.ApplyFunc(job.New, func(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR,
		body job.Body, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return expJob, nil
	})

	j, err := p.CreateJob("", "", 1, 1, "", true)
	require.Equal(t, expErr, err)
	require.Nil(t, j)
}

func TestCreateJobUnlockErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expErr := errors.New("lock error")
	testCases := []struct {
		unlockResult bool
		unlockErr    error
		isError      bool
	}{
		{unlockResult: false, unlockErr: expErr, isError: true},
		{unlockResult: true, unlockErr: nil, isError: false},
		{unlockResult: false, unlockErr: nil, isError: true},
	}

	for i, testCase := range testCases {
		t.Logf("run test case %d", i)
		mockStore := store.NewMockStore(ctrl)
		mockStore.EXPECT().CreateJob(gomock.Any()).Return(nil)

		mockLogger := log.NewMockLogger(ctrl)
		mockLogger.EXPECT().WithModule(gomock.Any()).Return(mockLogger)
		if testCase.isError {
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
		}

		mockLock := mocklock.NewMockLocker(ctrl)
		expJob := &job.Job{Mutex: mockLock}
		mockLock.EXPECT().Lock().Return(nil)
		mockLock.EXPECT().Unlock().Return(testCase.unlockResult, testCase.unlockErr)
		p := New(mockStore, mockLogger)

		gomonkey.ApplyFunc(job.New, func(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR,
			body job.Body, lockerFunc lock.LockerFunc) (*job.Job, error) {
			return expJob, nil
		})

		j, err := p.CreateJob("", "", 1, 1, "", false)
		require.Nil(t, err)
		require.Equal(t, expJob, j)
	}
}

func TestCreateJobCreateOrReplace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expErr := errors.New("test error")
	testCases := []struct {
		isReplace bool
		expError  error
	}{
		{isReplace: false, expError: nil},
		{isReplace: false, expError: expErr},
		{isReplace: true, expError: nil},
		{isReplace: true, expError: expErr},
	}

	for i, testCase := range testCases {
		t.Logf("run test case %d", i)
		mockStore := store.NewMockStore(ctrl)
		if testCase.isReplace {
			mockStore.EXPECT().ReplaceJob(gomock.Any()).Return(testCase.expError)
		} else {
			mockStore.EXPECT().CreateJob(gomock.Any()).Return(testCase.expError)
		}

		mockLogger := log.NewMockLogger(ctrl)
		mockLogger.EXPECT().WithModule(gomock.Any()).Return(mockLogger)

		mockLock := mocklock.NewMockLocker(ctrl)
		expJob := &job.Job{Mutex: mockLock}
		mockLock.EXPECT().Lock().Return(nil)
		mockLock.EXPECT().Unlock().Return(true, nil)
		p := New(mockStore, mockLogger)

		gomonkey.ApplyFunc(job.New, func(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR,
			body job.Body, lockerFunc lock.LockerFunc) (*job.Job, error) {
			return expJob, nil
		})

		j, err := p.CreateJob("", "", 1, 1, "", testCase.isReplace)
		require.Equal(t, testCase.expError, err)
		if testCase.expError != nil {
			require.Nil(t, j)
		} else {
			require.Equal(t, expJob, j)
		}
	}
}

func TestLoadReadyJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	mockLogger := log.NewMockLogger(ctrl)
	mockLogger.EXPECT().WithModule(gomock.Any()).Return(mockLogger)
	ver := job.NewVersion()
	j := &job.Job{Version: ver}

	// test case 1 job.Get error
	jobErr := errors.New("job err")
	gomonkey.ApplyFunc(job.Get, func(topic job.Topic, id job.Id, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return nil, jobErr
	})
	p := New(mockStore, mockLogger)
	jRet, err := p.LoadReadyJob("", "", ver)
	require.Nil(t, jRet)
	require.Equal(t, jobErr, err)

	// test case 2: s.LoadJob error
	loadErr := errors.New("load err")
	gomonkey.ApplyFunc(job.Get, func(topic job.Topic, id job.Id, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return j, nil
	})
	mockStore.EXPECT().LoadJob(j).Return(loadErr)
	jRet, err = p.LoadReadyJob("", "", ver)
	require.Nil(t, jRet)
	require.Equal(t, loadErr, err)

	// test case 3 j.IsVersionSame error
	gomonkey.ApplyFunc(job.Get, func(topic job.Topic, id job.Id, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return j, nil
	})
	mockStore.EXPECT().LoadJob(j).Return(nil)
	jRet, err = p.LoadReadyJob("", "", job.NewVersion())
	require.Nil(t, jRet)
	require.Equal(t, ErrVersionNotSame, err)

	// test case 4 pass
	gomonkey.ApplyFunc(job.Get, func(topic job.Topic, id job.Id, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return j, nil
	})
	mockStore.EXPECT().LoadJob(j).Return(nil)
	jRet, err = p.LoadReadyJob("", "", ver)
	require.Equal(t, j, jRet)
	require.Nil(t, err)
}

func TestLoadDeleteJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockStore(ctrl)
	mockLogger := log.NewMockLogger(ctrl)
	mockLogger.EXPECT().WithModule(gomock.Any()).Return(mockLogger)
	mockLock := mocklock.NewMockLocker(ctrl)

	j := &job.Job{Mutex: mockLock}
	p := New(mockStore, mockLogger)

	// test case 1: job.Get
	jobGetErr := errors.New("job err")
	gomonkey.ApplyFunc(job.Get, func(topic job.Topic, id job.Id, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return nil, jobGetErr
	})
	err := p.DeleteJob("", "")
	require.Equal(t, jobGetErr, err)

	// test case 2: j.Lock error
	lockErr := errors.New("lock err")
	gomonkey.ApplyFunc(job.Get, func(topic job.Topic, id job.Id, lockerFunc lock.LockerFunc) (*job.Job, error) {
		return j, nil
	})
	mockLock.EXPECT().Lock().Return(lockErr)
	err = p.DeleteJob("", "")
	require.Equal(t, lockErr, err)

	//  TODO: test case 3: DeleteJob error and unlock error
}
