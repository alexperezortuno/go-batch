package batch

import "time"

type Config struct {
	Size        int
	Timeout     time.Duration
	Concurrency int
	Ordered     bool
	Buffer      int

	MaxRetries int
	Backoff    time.Duration
}

type Option func(*Config)

func WithSize(size int) Option {
	return func(c *Config) {
		c.Size = size
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

func WithConcurrency(n int) Option {
	return func(c *Config) {
		c.Concurrency = n
	}
}

func WithOrdered() Option {
	return func(c *Config) {
		c.Ordered = true
	}
}

func WithRetry(maxRetries int, backoff time.Duration) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
		c.Backoff = backoff
	}
}
