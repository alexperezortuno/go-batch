package batch

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

func isRetryable(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "40001", // serialization_failure
			"40P01": // deadlock_detected
			return true
		}
	}
	return false
}

func retry(
	ctx context.Context,
	maxAttempts int,
	baseDelay time.Duration,
	fn func() error,
) error {

	var err error

	for attempt := 0; attempt <= maxAttempts; attempt++ {

		err = fn()
		if err == nil {
			return nil
		}

		if err != nil && !isRetryable(err) {
			return err
		}

		// last intent → out
		if attempt == maxAttempts {
			break
		}

		// exponential backoff: base * 2^attempt
		delay := baseDelay * time.Duration(1<<attempt)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return err
}
