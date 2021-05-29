package integration

import (
	"fmt"
	"github.com/changsongl/delay-queue-client/client"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

// TODO: All testing in this file will be improved in the future.

// This is an integration test for delay queue.
// It will test add job, consume and remove.
func TestDelayQueueAddAndRemove(t *testing.T) {
	// push n jobs with delay within 1 min
	DelayTimeSeconds := 60
	Jobs := 1
	rand.Seed(time.Now().Unix())

	cli := client.NewClient(DelayQueueAddr)
	t.Logf("Running test for %d jobs", Jobs)

	t.Log("Adding test")
	for i := 0; i < Jobs; i++ {
		delayTime := rand.Intn(DelayTimeSeconds)
		j, err := client.NewJob("test-topic", fmt.Sprintf("test-%d", i), client.JobDelayOption(time.Duration(delayTime)))
		assert.NoError(t, err)

		err = cli.AddJob(j)
		assert.NoError(t, err)
	}

	// consume jobs

	// check after 1.5 min, all jobs should be done
}
