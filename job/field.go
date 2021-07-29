package job

import (
	"time"
)

// Topic job topic
type Topic string

// ID job ID
type ID string

// Delay job delay time
type Delay time.Duration

// TTR job time to run
type TTR time.Duration

// Body job body
type Body string

// Empty interface
type Empty interface {
	IsEmpty() bool
}

// IsEmpty topic is empty
func (t Topic) IsEmpty() bool {
	return t == ""
}

// IsEmpty id is empty
func (id ID) IsEmpty() bool {
	return id == ""
}

// IsEmpty body is empty
func (b Body) IsEmpty() bool {
	return b == ""
}

// IsEmpty delay is empty
func (d Delay) IsEmpty() bool {
	return d == 0
}

// IsEmpty ttr is empty
func (t TTR) IsEmpty() bool {
	return t == 0
}
