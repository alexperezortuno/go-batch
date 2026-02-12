package batch

import (
	"context"
	"sync"
	"time"
)

func (b *Batcher[T]) Start(
	ctx context.Context,
	processor Processor[T],
) {
	b.processor = processor

	batches := make(chan []T, b.cfg.Concurrency)

	go b.runBatcher(ctx, batches)
	go b.runWorkers(ctx, batches)
}

func (b *Batcher[T]) runBatcher(
	ctx context.Context,
	out chan []T,
) {
	var ticker *time.Ticker

	if b.cfg.Timeout > 0 {
		ticker = time.NewTicker(b.cfg.Timeout)
		defer ticker.Stop()
	}

	batch := make([]T, 0, b.cfg.Size)

	flush := func() {
		if len(batch) > 0 {
			out <- batch
			batch = make([]T, 0, b.cfg.Size)
		}
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			close(out)
			return

		case item := <-b.in:
			batch = append(batch, item)
			if len(batch) >= b.cfg.Size {
				flush()
			}

		case <-tickerChan(ticker):
			flush()
		}
	}
}

func tickerChan(t *time.Ticker) <-chan time.Time {
	if t == nil {
		return nil
	}
	return t.C
}

func (b *Batcher[T]) runWorkers(
	ctx context.Context,
	in chan []T,
) {
	var wg sync.WaitGroup

	for i := 0; i < b.cfg.Concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for batch := range in {
				_ = b.processor(ctx, batch)
			}
		}()
	}

	wg.Wait()
}
