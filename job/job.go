package job

import (
	"errors"
	"fmt"
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
}

func New(topic Topic, id Id, delay Delay, ttr TTR, body Body) (*Job, error) {
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

	j.generateFingerprint()

	return j, nil
}

func (j *Job) getName() string {
	return fmt.Sprintf("%s_%s", j.topic, j.id)
}

// TODO: Lock
func (j *Job) Lock() error {
	return nil
}

// TODO: Unlock
func (j *Job) Unlock() error {
	return nil
}

// TODO: generateFingerprint
func (j *Job) generateFingerprint() {
	j.fp = 1
}
