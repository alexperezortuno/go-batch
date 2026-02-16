package batch

type orderer struct {
	next   int64
	buffer map[int64]batchResult
}

func newOrderer() *orderer {
	return &orderer{
		next:   1,
		buffer: make(map[int64]batchResult),
	}
}

func (o *orderer) handle(res batchResult) []batchResult {
	o.buffer[res.seq] = res

	var ready []batchResult

	for {
		r, ok := o.buffer[o.next]
		if !ok {
			break
		}

		ready = append(ready, r)
		delete(o.buffer, o.next)
		o.next++
	}

	return ready
}
