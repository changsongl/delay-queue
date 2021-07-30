package encode

import (
	"testing"
)

func TestCompressEncoder(t *testing.T) {
	encoder := NewCompress()
	runEncodeTest(t, encoder)
}
