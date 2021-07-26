package json

import (
	"encoding/json"
	"github.com/changsongl/delay-queue/config/decode"
)

// decoder for json
type decoder struct {
	dcFunc decode.Func
}

// NewDecoder new decoder
func NewDecoder() decode.Decoder {
	return &decoder{
		dcFunc: json.Unmarshal,
	}
}

// DecodeFunc decode function
func (d *decoder) DecodeFunc() decode.Func {
	return d.dcFunc
}
