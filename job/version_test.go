package job

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestNewVersionAndString(t *testing.T) {
	before := time.Now().UnixNano()
	ver := NewVersion()
	after := time.Now().UnixNano()
	require.LessOrEqual(t, before, ver.t.UnixNano(), "before should be less than version time")
	require.LessOrEqual(t, ver.t.UnixNano(), after, "version time should be less than after")

	verTime, err := strconv.ParseInt(ver.String(), 10, 64)
	require.NoError(t, err)
	require.LessOrEqual(t, before, verTime, "before should be less than version time")
	require.LessOrEqual(t, verTime, after, "version time should be less than after")
}

func TestNewLoadVersionAndEqual(t *testing.T) {
	ver := NewVersion()
	verNew, err := LoadVersion(ver.String())
	require.NoError(t, err)
	equal := ver.Equal(verNew)
	require.Equal(t, equal, true, "version should be same")
}

func TestVersionMarshalAndUnMarshall(t *testing.T) {
	ver := NewVersion()
	bytes, err := ver.MarshalJSON()
	require.NoError(t, err)

	verSame := Version{}
	err = verSame.UnmarshalJSON(bytes)
	require.NoError(t, err)

	equal := ver.Equal(verSame)
	require.Equal(t, equal, true, "version should be same")
}
