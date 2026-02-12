package batch

type batchWithSeq[T any] struct {
	seq   int64
	items []T
}

type batchResult struct {
	seq int64
	err error
}
