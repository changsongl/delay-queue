package encode

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncoder(t *testing.T) {
	encoder := New()
	s := struct {
		A string
		B string
	}{A: "A", B: "B"}

	str, err := encoder.Encode(s)
	require.NoError(t, err)

	sCopy := s
	sCopy.A = "AAA"
	sCopy.B = "BBB"
	err = encoder.Decode(str, &sCopy)
	require.NoError(t, err)
	require.EqualValues(t, s, sCopy)
}
