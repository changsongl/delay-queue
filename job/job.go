package job

import (
	"errors"
	"fmt"
	"time"

	"github.com/changsongl/delay-queue/pkg/lock"
)

type Job struct {
	Topic   Topic       `json:"topic"`
	ID      Id          `json:"id"`
	Delay   Delay       `json:"delay"`
	TTR     TTR         `json:"ttr"`
	Body    Body        `json:"body"`
	Version Version     `json:"version"`
	Mutex   lock.Locker `json:"-"`
}

// New return a job with everything init
func New(topic Topic, id Id, delay Delay, ttr TTR, body Body, lockerFunc lock.LockerFunc) (*Job, error) {
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
func Get(topic Topic, id Id, lockerFunc lock.LockerFunc) (*Job, error) {
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

func (j *Job) GetNameWithVersion() NameVersion {
	return NameVersion(fmt.Sprintf("%s_%s_%s", j.Topic, j.ID, j.Version))
}

// GetName return job lock name
func (j *Job) getLockName() string {
	return fmt.Sprintf("%s_lock", j.GetName())
}

// GetName return job unique name getter
func (j *Job) GetDelayTimeFromNow() time.Time {
	return time.Now().Add(time.Duration(j.Delay))
}

// GetName return job unique name getter
func (j *Job) GetTTRTimeFromNow() time.Time {
	return time.Now().Add(time.Duration(j.Delay))
}

// Lock lock the job
func (j *Job) Lock() error {
	return j.Mutex.Lock()
}

// Unlock unlock the job
func (j *Job) Unlock() (bool, error) {
	return j.Mutex.Unlock()
}
