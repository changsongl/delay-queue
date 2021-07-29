package job

import (
	"strconv"
	"time"
)

// Version job version, it is a time object. nano timestamp
type Version struct {
	t time.Time
}

// NewVersion create a new version
func NewVersion() Version {
	return Version{t: time.Now()}
}

// String function
func (v Version) String() string {
	return strconv.FormatInt(v.t.UnixNano(), 10)
}

// UInt64 function
func (v Version) UInt64() uint64 {
	return uint64(v.t.UnixNano())
}

// LoadVersion load version from a string
func LoadVersion(vs string) (Version, error) {
	version := Version{}
	vi, err := strconv.ParseInt(vs, 10, 64)
	if err != nil {
		return version, err
	}

	return Version{t: time.Unix(vi/1e9, vi%1e9)}, nil
}

// Equal check version v equals to version v2
func (v Version) Equal(v2 Version) bool {
	return v.t.UnixNano() == v2.t.UnixNano()
}

// MarshalJSON json marshall
func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalJSON json unmarshall
func (v *Version) UnmarshalJSON(b []byte) error {
	vi, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	v.t = time.Unix(vi/1e9, vi%1e9)
	return nil
}
