package job

import (
	"errors"
	"fmt"
	"strings"
)

type NameVersion string

func NewNameVersionString(str string) NameVersion {
	return NameVersion(str)
}

func NewNameVersion(topic Topic, id Id, version Version) NameVersion {
	return NameVersion(fmt.Sprintf("%s_%s_%s", topic, id, version))
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
