package job

import (
	"errors"
	"fmt"
	"strings"
)

// NameVersion name version, combination of topic, id and version
type NameVersion string

// NewNameVersionString convert string to name version
func NewNameVersionString(str string) NameVersion {
	return NameVersion(str)
}

// NewNameVersion create new name version object
func NewNameVersion(topic Topic, id ID, version Version) NameVersion {
	return NameVersion(fmt.Sprintf("%s_%s_%s", topic, id, version))
}

// MarshalBinary for json encode
func (nv NameVersion) MarshalBinary() ([]byte, error) {
	return []byte(nv), nil
}

// Parse parse name version to topic, id and version
func (nv NameVersion) Parse() (Topic, ID, Version, error) {
	data := strings.Split(string(nv), "_")
	if len(data) != 3 {
		return "", "", Version{}, errors.New("name version parse failed")
	}
	topic, id, vStr := Topic(data[0]), ID(data[1]), data[2]
	v, err := LoadVersion(vStr)
	if err != nil {
		return "", "", Version{}, err
	}

	return topic, id, v, nil
}
