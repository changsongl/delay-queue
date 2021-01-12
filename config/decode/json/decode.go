package json

import (
	"encoding/json"
	"github.com/changsongl/delay-queue/config/decode"
)

// decoder for json
type decoder struct {
	dcFunc decode.Func
}

func NewDecoder() decode.Decoder {
	return &decoder{
		dcFunc: json.Unmarshal,
	}
}

func (d *decoder) DecodeFunc() decode.Func {
	return d.dcFunc
}
