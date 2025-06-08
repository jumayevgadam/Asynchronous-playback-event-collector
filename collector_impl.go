package homework_5_1

import (
	"context"
	"fmt"
	"homework_5_1/model"
	"homework_5_1/stats"
	"sync"
	"sync/atomic"
	"time"
)

const (
	fullProgressPercent = 100 // 100% progress updates completion.
	maxProgressChan     = 100 // Buffer size of progress channel.
)

func init() {
	CurrentCollector = newCollectorImpl()
}

// collectorImpl manages event processing, enrichment and stats tracking.
type collectorImpl struct {
	enrichClient     EnrichClient // for fetching region and model data.
	totalEvents      uint64       // for tracking the total number of events processed.
	startTime        time.Time    // for calculate ProcessedPerSec.
	mu               sync.Mutex   // for protect shared state.
	enrichClientOnce sync.Once    // for ensuring enrichClient is set only once.
}

// newCollectorImpl creates a new collectorImpl with initialized start time.
func newCollectorImpl() *collectorImpl {
	return &collectorImpl{
		startTime: time.Now(),
	}
}

// operationImpl represents a single event processing operation.
type operationImpl struct {
	progressChan chan uint8       // channel for sending progress updates (0-100%).
	doneChan     chan struct{}    // channel signaling operation completion.
	result       *OperationResult // result of the operation.
	startTime    time.Time        // start time for this operation.
	processed    uint64           // number of events processed in this operation.
	mu           sync.Mutex       // protects specific operation state.
}

// newOperation creates a new operationImpl.
func newOperationImpl(totalEvents int) *operationImpl {
	op := &operationImpl{
		progressChan: make(chan uint8, maxProgressChan),
		doneChan:     make(chan struct{}),
		result: &OperationResult{
			HandledEvents: 0,
			Elapsed:       0,
		},
		startTime: time.Now(),
	}
	// handle empty events case.
	if totalEvents == 0 {
		op.progressChan <- fullProgressPercent
		close(op.progressChan)
		close(op.doneChan)
	}

	return op
}

// WithEnrichClient sets the EnrichClient only once and returns the collector.
func (c *collectorImpl) WithEnrichClient(client EnrichClient) EventCollector {
	c.enrichClientOnce.Do(func() {
		c.mu.Lock()         // we need to take guarantee the client is set
		defer c.mu.Unlock() // exactly once even if called multiple goroutines.
		c.enrichClient = client
	})

	return c
}

func (c *collectorImpl) Handle(events []*model.PlaybackEvent) (Operation, error) {
	if c.enrichClient == nil {
		return nil, fmt.Errorf("enrich client not set")
	}
	// initialize the operation to track progress and results.
	op := newOperationImpl(len(events))
	if len(events) == 0 {
		return op, nil
	}

	var wg sync.WaitGroup
	ctx := context.Background()
	// launch a goroutine to process all events concurrently.
	go func() {
		defer close(op.progressChan) // close progress channel when done.
		defer close(op.doneChan)     // signal operation completion.

		for _, event := range events {
			wg.Add(1)
			// Start a goroutine to process each event independently.
			go func(e *model.PlaybackEvent) {
				defer wg.Done()

				// Enrich event with region data.
				if region, _, err := c.enrichClient.GetRegion(ctx, e.UserID); err == nil {
					e.Region = region
				}

				// Enrich event with model data.
				if deviceModel, _, err := c.enrichClient.GetModel(ctx, e.DeviceType); err == nil {
					e.Model = deviceModel // fixed name.
				}

				// Safely update counters.
				newProcessed := atomic.AddUint64(&op.processed, 1)
				atomic.AddUint64(&c.totalEvents, 1)

				// Safely update result.
				op.mu.Lock()
				op.result.HandledEvents = newProcessed
				op.mu.Unlock()

				// calculate and send progress percentage.
				progress := uint8((float64(newProcessed) / float64(len(events))) * float64(fullProgressPercent))
				select {
				case op.progressChan <- progress:
				default:
				}
			}(event)
		}

		wg.Wait()

		// ensure final fullProgressPercent is marked as 100%.
		select {
		case op.progressChan <- fullProgressPercent:
		default:
		}

		// record total elapsed time for the operation.
		op.mu.Lock()
		op.result.Elapsed = time.Since(op.startTime)
		op.mu.Unlock()
	}()

	return op, nil
}

// Stats returns global statistics for the collector, including total events and processing rate.
func (c *collectorImpl) Stats() (*stats.Stats, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elapsed := time.Since(c.startTime).Seconds()
	var processedPerSec float64

	if elapsed > 0 {
		processedPerSec = float64(c.totalEvents) / elapsed
	}

	return &stats.Stats{
		TotalEvents:     c.totalEvents,
		ProcessedPerSec: processedPerSec,
	}, nil
}

// ProgressChan returns a read-only channel for receiving progress updates.
func (o *operationImpl) ProgressChan() <-chan uint8 {
	return o.progressChan
}

// Done returns a channel that closes when the operation is complete.
func (o *operationImpl) Done() chan struct{} {
	return o.doneChan
}

// Stats returns a snapshot of operation-specific statistics.
func (o *operationImpl) Stats() (*stats.Stats, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	elapsed := time.Since(o.startTime).Seconds()
	var processedPerSec float64

	if elapsed > 0 {
		processedPerSec = float64(o.processed) / elapsed
	}

	return &stats.Stats{
		TotalEvents:     o.processed,
		ProcessedPerSec: processedPerSec,
	}, nil
}

// Result returns the operation's result, including handled events and elapsed time.
func (o *operationImpl) Result() *OperationResult {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.result.Elapsed = time.Since(o.startTime)

	return o.result
}
