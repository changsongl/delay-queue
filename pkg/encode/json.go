package encode

import (
	"github.com/changsongl/delay-queue/job"
	jsoniter "github.com/json-iterator/go"
)

// implemented Encoder interface
type jsonEncoder struct {
}

// NewJSON create a json encoder
func NewJSON() Encoder {
	return &jsonEncoder{}
}

// Encode encode function
func (j jsonEncoder) Encode(i *job.Job) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Decode decode function
func (j jsonEncoder) Decode(b []byte, i *job.Job) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	return json.Unmarshal(b, i)
}
