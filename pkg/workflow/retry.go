package workflow

import (
	"context"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

func WithRetry(policy RetryPolicy, do func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		var lastErr error
		delay := policy.Delay

		for i := 1; i <= policy.MaxAttempts; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if err := do(ctx); err == nil {
				return nil
			} else {
				lastErr = err
			}

			if i < policy.MaxAttempts {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
					if policy.Multiplier > 1 {
						delay = time.Duration(float64(delay) * policy.Multiplier)
					}
				}
			}
		}
		return lastErr
	}
}
