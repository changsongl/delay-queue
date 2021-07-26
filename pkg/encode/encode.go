package encode

import "encoding/json"

// Encoder encoder interface
type Encoder interface {
	Encode(interface{}) (string, error)
	Decode(string, interface{}) error
}

// implemented Encoder interface
type jsonEncoder struct {
}

// New create an encoder
func New() Encoder {
	return &jsonEncoder{}
}

// Encode encode function
func (j jsonEncoder) Encode(i interface{}) (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Decode decode function
func (j jsonEncoder) Decode(s string, i interface{}) error {
	return json.Unmarshal([]byte(s), i)
}
