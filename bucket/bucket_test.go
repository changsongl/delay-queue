package bucket

import (
	"errors"
	"github.com/changsongl/delay-queue/job"
	mock_log "github.com/changsongl/delay-queue/test/mock/log"
	lockmock "github.com/changsongl/delay-queue/test/mock/pkg/lock"
	storemock "github.com/changsongl/delay-queue/test/mock/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBucketCreateJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := storemock.NewMockStore(ctrl)
	lockMk := lockmock.NewMockLocker(ctrl)
	mLog := mock_log.NewMockLogger(ctrl)

	sm.EXPECT().GetLock(gomock.All()).Return(lockMk).AnyTimes()
	b := New(sm, mLog, 10, "test_bucket")

	// case1: no error
	sm.EXPECT().CreateJobInBucket(gomock.Eq("test_bucket_1"), gomock.All(), gomock.All()).Return(nil)
	err := b.CreateJob(nil, true)
	require.NoError(t, err, "first create should no error")

	// case2: has error
	expectErr := errors.New("expect error")
	sm.EXPECT().CreateJobInBucket(gomock.Eq("test_bucket_2"), gomock.All(), gomock.All()).Return(expectErr)
	err = b.CreateJob(nil, true)
	require.Equal(t, expectErr, err, "second create should be expect error")
}

func TestBucketGetBuckets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := storemock.NewMockStore(ctrl)
	lockMk := lockmock.NewMockLocker(ctrl)
	mLog := mock_log.NewMockLogger(ctrl)

	sm.EXPECT().GetLock(gomock.All()).Return(lockMk).AnyTimes()
	b := New(sm, mLog, 2, "test_bucket")
	bucketNames := b.GetBuckets()

	expectNames := []uint64{
		0, 1,
	}

	for i, bucketName := range bucketNames {
		if i > len(expectNames) {
			t.Error("it is greater than expecting length")
			t.FailNow()
		}

		require.Equal(t, expectNames[i], bucketName, "bucket names are not equal")
	}
}

func TestBucketGetBucketJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := storemock.NewMockStore(ctrl)
	lockMk := lockmock.NewMockLocker(ctrl)
	mLog := mock_log.NewMockLogger(ctrl)

	sm.EXPECT().GetLock(gomock.All()).Return(lockMk).AnyTimes()
	b := New(sm, mLog, 2, "test_bucket")

	expectErr := errors.New("error GetReadyJobsInBucket")
	sm.EXPECT().GetReadyJobsInBucket(gomock.Eq("test_bucket_0"), gomock.All()).Return(nil, expectErr)
	versions, err := b.GetBucketJobs(0)
	require.Equal(t, expectErr, err, "it should be expecting")
	require.Nil(t, versions, "version names should be nil")

	expectNvs := []job.NameVersion{
		"nv1", "nv2",
	}
	sm.EXPECT().GetReadyJobsInBucket(gomock.Eq("test_bucket_1"), gomock.All()).Return(expectNvs, nil)
	versions, err = b.GetBucketJobs(1)
	require.NoError(t, err, "it should have no error")
	require.Equal(t, expectNvs, versions, "version names should be equal")
}

func TestBucketFetchNum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := storemock.NewMockStore(ctrl)
	lockMk := lockmock.NewMockLocker(ctrl)
	mLog := mock_log.NewMockLogger(ctrl)

	sm.EXPECT().GetLock(gomock.All()).Return(lockMk).AnyTimes()
	b := New(sm, mLog, 2, "test_bucket")
	require.Equal(t, DefaultMaxFetchNum, b.GetMaxFetchNum(), "fetch number should be default")

	var newNum uint64 = 30
	b.SetMaxFetchNum(newNum)
	require.Equal(t, newNum, b.GetMaxFetchNum(), "fetch number should be new number")
}

// TODO: test collect metric
