package workflow

import (
	"context"
	"time"

	"github.com/mirkobrombin/go-foundation/pkg/resiliency"
)

type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

func WithRetry(policy RetryPolicy, do func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		return resiliency.Retry(ctx, func() error {
			return do(ctx)
		},
			resiliency.WithAttempts(policy.MaxAttempts),
			resiliency.WithDelay(policy.Delay, 24*time.Hour),
			resiliency.WithFactor(policy.Multiplier),
		)
	}
}
