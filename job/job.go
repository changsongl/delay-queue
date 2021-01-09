package job

import (
	"errors"
	"fmt"
	"github.com/changsongl/delay-queue/pkg/lock"
	"time"
)

type Job struct {
	Topic Topic       `json:"topic"`
	ID    Id          `json:"id"`
	Delay Delay       `json:"delay"`
	TTR   TTR         `json:"ttr"`
	Body  Body        `json:"body"`
	TS    Version     `json:"version"`
	Mutex lock.Locker `json:"-"`
}

// New return a job with everything init
func New(topic Topic, id Id, delay Delay, ttr TTR, body Body, lockerFunc lock.LockerFunc) (*Job, error) {
	if topic == "" || id == "" {
		return nil, errors.New("[job.New] Topic or ID is empty")
	}

	j := &Job{
		Topic: topic,
		ID:    id,
		Delay: delay,
		TTR:   ttr,
		Body:  body,
		TS:    NewVersion(),
	}

	j.Mutex = lockerFunc(j.getLockName())

	return j, nil
}

// GetName return job unique name getter
func (j *Job) GetName() string {
	return fmt.Sprintf("%s_%s", j.Topic, j.ID)
}

func (j *Job) GetNameWithVersion() NameVersion {
	return NameVersion(fmt.Sprintf("%s_%s_%v", j.Topic, j.ID, j.TS))
}

// GetName return job lock name
func (j *Job) getLockName() string {
	return fmt.Sprintf("%s_lock", j.GetName())
}

// GetName return job unique name getter
func (j *Job) GetDelayTimeFromNow() time.Time {
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
