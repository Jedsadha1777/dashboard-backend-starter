package utils

import (
	"context"
	"time"
)

type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	RandomizeFactor float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:     3,
		InitialDelay:    100 * time.Millisecond,
		MaxDelay:        10 * time.Second,
		Multiplier:      2.0,
		RandomizeFactor: 0.1,
	}
}

func RetryWithBackoff(ctx context.Context, config RetryConfig, fn func() error) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err

			// Check if context is done
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Last attempt, don't sleep
			if attempt == config.MaxAttempts-1 {
				break
			}

			// Calculate next delay with exponential backoff
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}

			// Update delay for next iteration
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		} else {
			return nil // Success
		}
	}

	return lastErr
}

func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		if i < attempts-1 {
			time.Sleep(delay)
		}
	}

	return err
}
