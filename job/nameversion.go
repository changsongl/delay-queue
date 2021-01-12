package job

import (
	"errors"
	"strings"
)

type NameVersion string

func NewNameVersion(str string) NameVersion {
	return NameVersion(str)
}

func (nv NameVersion) MarshalBinary() ([]byte, error) {
	return []byte(nv), nil
}

func (nv NameVersion) Parse() (Topic, Id, Version, error) {
	data := strings.Split(string(nv), "_")
	if len(data) != 3 {
		return "", "", Version{}, errors.New("name version parse failed")
	}
	topic, id, vStr := Topic(data[0]), Id(data[1]), data[2]
	v, err := LoadVersion(vStr)
	if err != nil {
		return "", "", Version{}, err
	}

	return topic, id, v, nil
}
