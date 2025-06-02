package homework_5_1

import (
	"testing"

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
