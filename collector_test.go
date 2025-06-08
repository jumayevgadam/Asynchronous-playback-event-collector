package homework_5_1

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestDepartment is a simple test that checks if basic functionality
// works well.
func TestDepartment(t *testing.T) {
	enrichClient := newEnrichClient(t)

	collector := CurrentCollector.WithEnrichClient(enrichClient)
	events := generateEvents(t)

	assert := require.New(t)

	operation, err := collector.Handle(events)
	assert.NoError(err)

	progress := <-operation.ProgressChan()
	assert.NotEmpty(progress)

	<-operation.Done()

	assert.EqualValues(len(events), operation.Result().HandledEvents)
}

func TestProgressReporting(t *testing.T) {
	enrichClient := newSlowEnrichClient(t)

	collector := CurrentCollector.WithEnrichClient(enrichClient)
	events := generateEvents(t)

	assert := require.New(t)

	operation, err := collector.Handle(events)
	assert.NoError(err)

	var (
		progressSnapshots []uint8
		timeout           = time.After(50 * time.Second)
	)

loop:
	for {
		select {
		case p, ok := <-operation.ProgressChan():
			if !ok {
				break loop
			}
			progressSnapshots = append(progressSnapshots, p)
		case <-timeout:
			t.Fatal("Timeout waiting for progress updates")
		}
	}

	<-operation.Done()
	assert.GreaterOrEqual(len(progressSnapshots), 10, "Should have received at least 10 progress update")

	last := progressSnapshots[len(progressSnapshots)-1]
	assert.Equal(uint8(100), last, "Final progress update should be 100%%")
}

func TestProcessPerformance(t *testing.T) {
	enrichClient := newSlowEnrichClient(t)

	collector := CurrentCollector.WithEnrichClient(enrichClient)
	events := generateEvents(t)

	assert := require.New(t)

	operation, err := collector.Handle(events)
	assert.NoError(err)

	startTime := time.Now()

	var (
		elapsedTime float64
		timeout     = time.After(20 * time.Second)
	)

LOOP:
	for {
		select {
		case <-operation.Done():
			elapsedTime = time.Since(startTime).Seconds()
			break LOOP
		case <-timeout:
			t.Fatal("Timeout waiting for progress updates")
		}
	}

	assert.LessOrEqual(elapsedTime, float64(15), fmt.Sprintf("Should process less than 15 seconds, but elapsed time %f seconds", elapsedTime))
}
