package batch

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
)

type Processor[T any] func(ctx context.Context, items []T) error

type Batcher[T any] struct {
	cfg       Config
	in        chan T
	processor Processor[T]

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	started atomic.Bool
	closed  atomic.Bool
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

func (b *Batcher[T]) Start(
	parent context.Context,
	processor Processor[T],
) {
	if b.started.Load() {
		panic("batcher already started")
	}

	b.started.Store(true)

	b.ctx, b.cancel = context.WithCancel(parent)
	b.processor = processor

	batches := make(chan batchWithSeq[T], b.cfg.Concurrency)
	results := make(chan batchResult, b.cfg.Concurrency)

	// 🟢 Batcher
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.runBatcher(b.ctx, batches)
	}()

	// 🟢 Workers (ahora usando runWorkers)
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.runWorkers(b.ctx, batches, results)
	}()

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		o := newOrderer()

		for res := range results {

			ready := o.handle(res)

			for _, r := range ready {
				if r.err != nil {
					log.Printf("batch seq=%d failed: %v", r.seq, r.err)
				}
			}
		}
	}()
}

func (b *Batcher[T]) Submit(item T) error {
	if b.closed.Load() {
		return ErrClosed
	}

	select {
	case b.in <- item:
		return nil
	case <-b.ctx.Done():
		return b.ctx.Err()
	}
}

func (b *Batcher[T]) Stop() {
	if !b.closed.CompareAndSwap(false, true) {
		return
	}

	b.cancel()
	close(b.in)
	b.wg.Wait()
}
