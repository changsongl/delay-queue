package bucket

import (
	"errors"
	lockmock "github.com/changsongl/delay-queue/test/mock/pkg/lock"
	storemock "github.com/changsongl/delay-queue/test/mock/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBucketCreateJob(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sm := storemock.NewMockStore(ctrl)
	lockMk := lockmock.NewMockLocker(ctrl)

	b := New(sm, 10, "test_bucket")

	// case1: no error
	sm.EXPECT().GetLock(gomock.All()).Return(lockMk).AnyTimes()
	sm.EXPECT().CreateJobInBucket(gomock.Eq("test_bucket_1"), gomock.All(), gomock.All()).Return(nil)
	err := b.CreateJob(nil, true)
	require.NoError(t, err, "first create should no error")

	// case2: has error
	expectErr := errors.New("expect error")
	sm.EXPECT().CreateJobInBucket(gomock.Eq("test_bucket_2"), gomock.All(), gomock.All()).Return(expectErr)
	err = b.CreateJob(nil, true)
	require.Equal(t, expectErr, err, "second create should be expect error")
}
