package encode

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/stretchr/testify/require"
	"testing"
)

// BenchmarkCompressEncode    	1000000000	         0.000011 ns/op   len(32)
func BenchmarkCompressEncode(t *testing.B) {
	encoder := NewCompress()
	j, err := job.New("jobTopic", "131223", 1, 1, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	bytes, _ := encoder.Encode(j)

	_ = encoder.Decode(bytes, j)
}

// BenchmarkJSONEncode    	1000000000	         0.000024 ns/op    len(92)
func BenchmarkJSONEncode(t *testing.B) {
	encoder := NewJSON()
	j, err := job.New("jobTopic", "131223", 1, 1, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	bytes, _ := encoder.Encode(j)
	_ = encoder.Decode(bytes, j)
}
