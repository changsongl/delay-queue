package job

import (
	"errors"
	"fmt"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/changsongl/delay-queue/store"
	"time"
)

type Fingerprint uint

type Job struct {
	topic Topic
	id    Id
	delay Delay
	ttr   TTR
	body  Body
	ts    time.Time
	fp    Fingerprint
	mutex lock.Locker
}

// New return a job with everything init
func New(topic Topic, id Id, delay Delay, ttr TTR, body Body, store store.Store) (*Job, error) {
	if topic == "" || id == "" {
		return nil, errors.New("[job.New] topic or id is empty")
	}

	j := &Job{
		topic: topic,
		id:    id,
		delay: delay,
		ttr:   ttr,
		body:  body,
		ts:    time.Now(),
	}

	j.mutex = store.GetLock(j.GetName())
	j.generateFingerprint()

	return j, nil
}

// GetName return job unique name getter
func (j *Job) GetName() string {
	return fmt.Sprintf("%s_%s", j.topic, j.id)
}

// GetName return job unique name getter
func (j *Job) GetDelayTimeFromNow() time.Time {
	return time.Now().Add(time.Duration(j.delay))
}

// Lock lock the job
func (j *Job) Lock() error {
	return j.mutex.Lock()
}

// Unlock unlock the job
func (j *Job) Unlock() (bool, error) {
	return j.mutex.Unlock()
}

// TODO: generateFingerprint
func (j *Job) generateFingerprint() {
	j.fp = 1
}
