package encode

import (
	"github.com/changsongl/delay-queue/job"
)

// Encoder encoder interface
type Encoder interface {
	Encode(*job.Job) ([]byte, error)
	Decode([]byte, *job.Job) error
}
