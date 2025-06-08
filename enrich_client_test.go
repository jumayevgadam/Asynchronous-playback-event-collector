package homework_5_1

import (
	"context"
	"fmt"
	"homework_5_1/model"
	"testing"
	"time"
)

const (
	defaultBitrateKbps   = 128
	defaultRegion        = "Lebap"
	defaultModel         = "iPhone16 Pro Max"
	defaultFetchDuration = 1 * time.Second
)

type simpleTestClient struct {
	tb testing.TB
}

func newEnrichClient(tb testing.TB) *simpleTestClient {
	return &simpleTestClient{
		tb: tb,
	}
}

func (s *simpleTestClient) GetRegion(ctx context.Context, userID string) (region string, fetchDuration time.Duration, err error) {
	s.tb.Logf("GetRegion called with params userID=%s", userID)
	return defaultRegion, defaultFetchDuration, nil
}

func (s *simpleTestClient) GetModel(ctx context.Context, deviceType string) (model string, fetchDuration time.Duration, err error) {
	s.tb.Logf("GetModel called with params deviceType=%s", deviceType)
	return defaultModel, defaultFetchDuration, nil
}

func generateEvents(tb testing.TB) []*model.PlaybackEvent {
	const generateEventsNum = 10
	out := make([]*model.PlaybackEvent, generateEventsNum)

	for i := range generateEventsNum {
		event := model.PlaybackEvent{
			ID:          1,
			UserID:      fmt.Sprintf("u_%d", i),
			VideoID:     fmt.Sprintf("v_%d", i),
			StartAt:     int64(i),
			StopAt:      int64(2 * i),
			BitrateKbps: defaultBitrateKbps,
			DeviceType:  "Android: Samsung-G973F",
			ErrorCode:   nil,
		}

		tb.Logf("generated event %#v", event)

		out[i] = &event
	}

	return out
}
