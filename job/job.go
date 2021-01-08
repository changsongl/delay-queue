package job

import "github.com/changsongl/delay-queue/type/job"

type Fingerprint uint

type Job struct {
	topic job.Topic
	id    job.Id
	delay job.Delay
	ttr   job.TTR
	body  job.Body
	fp    Fingerprint
}

func New(topic job.Topic, id job.Id, delay job.Delay, ttr job.TTR, body job.Body) *Job {
	j := &Job{
		topic: topic,
		id:    id,
		delay: delay,
		ttr:   ttr,
		body:  body,
	}
	j.generateFingerprint()

	return j
}

func (j *Job) generateFingerprint() {
	j.fp = 1
}
