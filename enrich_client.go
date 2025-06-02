package homework_5_1

import (
	"context"
	"time"
)

type EnrichClient interface {
	GetRegion(ctx context.Context, userID string) (string, time.Duration, error)
	GetModel(ctx context.Context, deviceType string) (string, time.Duration, error)
}
