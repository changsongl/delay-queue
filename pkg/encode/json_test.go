package encode

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncoder(t *testing.T) {
	encoder := NewJSON()
	j, err := job.New("jobTopic", "1", 1, 1, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	str, err := encoder.Encode(j)
	require.NoError(t, err)

	jDecode := &job.Job{}
	err = encoder.Decode(str, jDecode)
	fmt.Println(j.Version.String())

	require.NoError(t, err)
	require.Equal(t, j.ID, jDecode.ID)
	require.Equal(t, j.TTR, jDecode.TTR)
	require.Equal(t, j.Delay, jDecode.Delay)
	require.Equal(t, j.Topic, jDecode.Topic)
	require.True(t, j.Version.Equal(jDecode.Version))
	require.Equal(t, j.Body, jDecode.Body)
	require.Equal(t, j.Version.String(), jDecode.Version.String())
}
