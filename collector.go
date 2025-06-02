package homework_5_1

import (
	"fmt"
	"homework_5_1/model"
	"homework_5_1/stats"
)

var CurrentCollector EventCollector = nil

type EventCollector interface {
	// WithEnrichClient sets enrich client.
	WithEnrichClient(client EnrichClient) EventCollector
	// Handle spawns an operation that serves events. Should not block
	// the execution.
	Handle(events []*model.PlaybackEvent) (Operation, error)
}

var errNotImpl = fmt.Errorf("not implemented")

type collectorNOOP struct {
}

var _ EventCollector = (*collectorNOOP)(nil)

func (c collectorNOOP) WithEnrichClient(client EnrichClient) EventCollector {
	return c
}

func (collectorNOOP) Handle(events []*model.PlaybackEvent) (Operation, error) {
	return nil, errNotImpl
}

func (collectorNOOP) Stats() (*stats.Stats, error) {
	return nil, errNotImpl
}
