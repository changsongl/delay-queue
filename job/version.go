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
	return strconv.FormatInt(v.t.UnixNano(), 10)
}

func LoadVersion(vs string) (Version, error) {
	version := Version{}
	vi, err := strconv.ParseInt(vs, 10, 64)
	if err != nil {
		return version, err
	}

	return Version{t: time.Unix(vi/1e9, vi%1e9)}, nil
}

func (v Version) Equal(v2 Version) bool {
	return v.t.UnixNano() == v2.t.UnixNano()
}

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) UnmarshalJSON(b []byte) error {
	vi, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	v.t = time.Unix(vi/1e9, vi%1e9)
	return nil
}
