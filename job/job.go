package job

import (
	"errors"
	"fmt"
	"time"

	"github.com/changsongl/delay-queue/pkg/lock"
)

// Job job for delay queue
type Job struct {
	Topic   Topic       `json:"topic,omitempty"`
	ID      ID          `json:"id,omitempty"`
	Delay   Delay       `json:"delay,omitempty"`
	TTR     TTR         `json:"ttr,omitempty"`
	Body    Body        `json:"body,omitempty"`
	Version Version     `json:"version,omitempty"`
	Mutex   lock.Locker `json:"-"`
}

// New return a job with everything init
func New(topic Topic, id ID, delay Delay, ttr TTR, body Body, lockerFunc lock.LockerFunc) (*Job, error) {
	j := &Job{
		Topic:   topic,
		ID:      id,
		Delay:   delay,
		TTR:     ttr,
		Body:    body,
		Version: NewVersion(),
	}

	err := j.IsValid()
	if err != nil {
		return nil, err
	}

	j.Mutex = lockerFunc(j.getLockName())

	return j, nil
}

// Get a job entity before load all information from storage
func Get(topic Topic, id ID, lockerFunc lock.LockerFunc) (*Job, error) {
	j := &Job{
		Topic: topic,
		ID:    id,
	}
	err := j.IsValid()
	if err != nil {
		return nil, err
	}

	j.Mutex = lockerFunc(j.getLockName())
	return j, nil
}

// IsValid check job is valid. job is not nil and topic and id
// is not empty
func (j *Job) IsValid() error {
	if j.Topic == "" || j.ID == "" {
		return errors.New("topic or id is empty")
	}

	return nil
}

// IsVersionSame return whether j's version is equal to v
func (j *Job) IsVersionSame(v Version) bool {
	return j.Version.Equal(v)
}

// GetName return job unique name getter
func (j *Job) GetName() string {
	return fmt.Sprintf("%s_%s", j.Topic, j.ID)
}

// GetNameWithVersion return name version of job
func (j *Job) GetNameWithVersion() NameVersion {
	return NewNameVersion(j.Topic, j.ID, j.Version)
}

// GetName return job lock name
func (j *Job) getLockName() string {
	return fmt.Sprintf("%s_lock", j.GetName())
}

// GetDelayTimeFromNow return how much time to wait for delaying
func (j *Job) GetDelayTimeFromNow() time.Time {
	return time.Now().Add(time.Duration(j.Delay))
}

// GetTTRTimeFromNow return how much time to wait until overtime
func (j *Job) GetTTRTimeFromNow() time.Time {
	return time.Now().Add(time.Duration(j.TTR))
}

// Lock lock the job
func (j *Job) Lock() error {
	return j.Mutex.Lock()
}

// Unlock unlock the job
func (j *Job) Unlock() (bool, error) {
	return j.Mutex.Unlock()
}

// SetVersion set job version by nano ts
func (j *Job) SetVersion(ts int64) {
	j.Version.t = time.Unix(ts/1e9, ts%1e9)
}
