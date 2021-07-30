package encode

import (
	"testing"
)

func TestEncoder(t *testing.T) {
	encoder := NewJSON()
	runEncodeTest(t, encoder)
}
