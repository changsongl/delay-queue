package integration

import (
	"fmt"
	"github.com/changsongl/delay-queue-client/client"
	"github.com/changsongl/delay-queue-client/consumer"
	"github.com/changsongl/delay-queue-client/job"
	"github.com/stretchr/testify/require"
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
	Jobs := 200
	topic, key := "test-topic", "testing_set"
	rand.Seed(time.Now().Unix())

	cli := client.NewClient(DelayQueueAddr)
	t.Logf("Running test for %d jobs", Jobs)

	t.Log("Adding test")
	for i := 0; i < Jobs; i++ {
		delayTime := rand.Intn(DelayTimeSeconds)
		id := fmt.Sprintf("test-%d", i)
		j, err := job.New(topic, id, job.JobDelayOption(time.Duration(delayTime)*time.Second))
		require.NoError(t, err)

		err = AddJobRecord(key, id)
		require.NoError(t, err)

		err = cli.AddJob(j)
		require.NoError(t, err)
	}

	t.Log("Finish adding and consume")

	go func() {
		m := make(map[string]int)
		// consume jobs
		c := consumer.New(cli, topic, consumer.WorkerNumOption(1))
		ch := c.Consume()
		for jobMsg := range ch {
			id := jobMsg.GetId()
			err := DeleteJobRecord(key, id)
			require.NoError(t, err)

			m[id]++
			if m[id] > 1 {
				t.Errorf("job id (%s) consume more than 1 time", id)
			}

			err = jobMsg.Finish()
			require.NoError(t, err)
		}
	}()

	// check after 1.5 min, all jobs should be done
	t.Log("Sleeping")
	time.Sleep(90 * time.Second)

	num, err := RecordNumbers(key)
	require.NoError(t, err)
	require.Equal(t, int64(0), num, "Remain jobs should be empty")
}
