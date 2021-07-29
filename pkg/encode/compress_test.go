package encode

import (
	"fmt"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompressEncoder(t *testing.T) {
	encoder := NewCompress()
	j, err := job.New("jobTopic11", "1222", 10, 5, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	str, err := encoder.Encode(j)
	require.NoError(t, err)
	fmt.Println(str)

	jDecode := &job.Job{}
	err = encoder.Decode(str, jDecode)
	fmt.Printf("%+v", jDecode)

	require.NoError(t, err)
	require.Equal(t, j.ID, jDecode.ID)
	require.Equal(t, j.TTR, jDecode.TTR)
	require.Equal(t, j.Delay, jDecode.Delay)
	require.Equal(t, j.Topic, jDecode.Topic)
	require.True(t, j.Version.Equal(jDecode.Version))
	require.Equal(t, j.Body, jDecode.Body)
	require.Equal(t, j.Version.String(), jDecode.Version.String())
}
