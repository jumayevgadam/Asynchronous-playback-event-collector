package homework_5_1

import (
	"context"
	"testing"
	"time"
)

type slowTestClient struct {
	tb testing.TB
}

func newSlowEnrichClient(tb testing.TB) *slowTestClient {
	return &slowTestClient{
		tb: tb,
	}
}

func (s *slowTestClient) GetRegion(ctx context.Context, userID string) (region string, fetchDuration time.Duration, err error) {
	s.tb.Logf("GetRegion called with params userID=%s", userID)
	time.Sleep(5 * time.Second)

	return defaultRegion, defaultFetchDuration, nil
}

func (s *slowTestClient) GetModel(ctx context.Context, deviceType string) (model string, fetchDuration time.Duration, err error) {
	time.Sleep(5 * time.Second)
	s.tb.Logf("GetModel called with params deviceType=%s", deviceType)

	return defaultModel, defaultFetchDuration, nil
}
