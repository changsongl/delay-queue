package job

import (
	"strconv"
	"time"
)

type Version struct {
	t time.Time
}

func NewVersion() Version {
	return Version{t: time.Now()}
}

func (v Version) String() string {
	return strconv.Itoa(v.t.Nanosecond())
}

func LoadVersion(vs string) (Version, error) {
	version := Version{}
	vi, err := strconv.ParseInt(vs, 10, 8)
	if err != nil {
		return version, err
	}

	return Version{t: time.Unix(vi/100000, vi%100000)}, nil
}

func (v Version) Equal(v2 Version) bool {
	return v.t.Nanosecond() == v2.t.Nanosecond()
}
