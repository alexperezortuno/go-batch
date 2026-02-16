package batch

import (
	"context"
	"sync"
	"time"
)

func (b *Batcher[T]) runBatcher(
	ctx context.Context,
	out chan batchWithSeq[T],
) {
	defer close(out)

	var ticker *time.Ticker
	if b.cfg.Timeout > 0 {
		ticker = time.NewTicker(b.cfg.Timeout)
		defer ticker.Stop()
	}

	batch := make([]T, 0, b.cfg.Size)
	var seq int64

	flush := func() {
		if len(batch) == 0 {
			return
		}

		seq++

		out <- batchWithSeq[T]{
			seq:   seq,
			items: batch,
		}

		batch = make([]T, 0, b.cfg.Size)
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return

		case item, ok := <-b.in:
			if !ok {
				flush()
				return
			}

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
	batches chan batchWithSeq[T],
	results chan batchResult,
) {
	var wg sync.WaitGroup

	for i := 0; i < b.cfg.Concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for batch := range batches {

				err := retry(
					ctx,
					b.cfg.MaxRetries,
					b.cfg.Backoff,
					func() error {
						return b.processor(ctx, batch.items)
					},
				)

				results <- batchResult{
					seq: batch.seq,
					err: err,
				}
			}
		}()
	}

	wg.Wait()
	close(results)
}

func (b *Batcher[T]) waitWorkersAndClose(
	results chan batchResult,
) {
	b.wg.Wait()
	close(results)
}
