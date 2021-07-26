package yaml

import (
	"github.com/changsongl/delay-queue/config/decode"
	"gopkg.in/yaml.v2"
)

// decoder for yaml
type decoder struct {
	dcFunc decode.Func
}

// NewDecoder new yaml decoder
func NewDecoder() decode.Decoder {
	return &decoder{
		dcFunc: yaml.Unmarshal,
	}
}

// DecodeFunc yaml decode function
func (d *decoder) DecodeFunc() decode.Func {
	return d.dcFunc
}
