package yaml

import (
	"github.com/changsongl/delay-queue/config/decode"
	"gopkg.in/yaml.v2"
)

type decoder struct {
	dcFunc decode.Func
}

func NewDecoder() decode.Decoder {
	return &decoder{
		dcFunc: yaml.Unmarshal,
	}
}

func (d *decoder) DecodeFunc() decode.Func {
	return d.dcFunc
}
