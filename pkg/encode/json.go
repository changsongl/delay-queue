package encode

import (
	"encoding/json"
	"github.com/changsongl/delay-queue/job"
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
	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Decode decode function
func (j jsonEncoder) Decode(b []byte, i *job.Job) error {
	return json.Unmarshal(b, i)
}
