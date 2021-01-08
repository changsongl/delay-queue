package encode

import "encoding/json"

type Encoder interface {
	Encode(interface{}) (string, error)
	Decode(string, interface{}) error
}

type jsonEncoder struct {
}

func New() Encoder {
	return &jsonEncoder{}
}

func (j jsonEncoder) Encode(i interface{}) (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (j jsonEncoder) Decode(s string, i interface{}) error {
	return json.Unmarshal([]byte(s), i)
}
