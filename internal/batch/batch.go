package batch

import (
	"context"
)

type Processor[T any] func(ctx context.Context, items []T) error

type Batcher[T any] struct {
	cfg       Config
	in        chan T
	processor Processor[T]
}

func New[T any](opts ...Option) *Batcher[T] {
	cfg := Config{
		Size:        100,
		Timeout:     0,
		Concurrency: 1,
		Buffer:      1000,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &Batcher[T]{
		cfg: cfg,
		in:  make(chan T, cfg.Buffer),
	}
}

func (b *Batcher[T]) Submit(item T) {
	b.in <- item
}

func (b *Batcher[T]) TrySubmit(item T) bool {
	select {
	case b.in <- item:
		return true
	default:
		return false
	}
}
