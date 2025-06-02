package homework_5_1

import (
	"homework_5_1/stats"
	"time"
)

type OperationResult struct {
	// HandledEvents is a total amount of events handled during this
	// operation
	HandledEvents uint64
	// Elapsed shows total time spent by this operation.
	Elapsed time.Duration
}

type Operation interface {
	// ProgressChan returns a channel for reporting current progress of handling
	// operation. It's expected to receive progress report each second.
	ProgressChan() <-chan uint8
	// Done report whether the operation
	Done() chan struct{}
	// Stats returns a snapshot of the current stats.
	Stats() (*stats.Stats, error)

	Result() *OperationResult
}
