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
