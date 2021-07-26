package job

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNameVersionMethods(t *testing.T) {
	topic := Topic("name")
	id := ID("id")
	version := NewVersion()

	nameVer := NewNameVersion(topic, id, version)
	nvBytes, err := nameVer.MarshalBinary()
	require.NoError(t, err)

	nameVerFromString := NewNameVersionString(string(nvBytes))
	require.Equal(t, nameVerFromString, nameVer, "name version should be same")

	topicParse, idParse, versionParse, err := nameVerFromString.Parse()
	require.NoError(t, err)
	require.Equal(t, topicParse, topic, "the topic of name version should be same")
	require.Equal(t, idParse, id, "the id of name version should be same")
	require.Equal(t, versionParse.String(), version.String(), "the version of name version should be same")

}
