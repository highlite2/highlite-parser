package cache

// IMemo is a memorization of a callback by key.
type IMemo interface {
	// Get takes the value for a key from cache, and if there is no such,
	// executes the callback and puts the result into the cache.
	Get(key string, callback MemoCallback) (interface{}, error)
}

// MemoCallback is a callback, that has to be executed to get a result for the key.
type MemoCallback func() (interface{}, error)

// memoResult is a result of a MemoCallback execution.
type memoResult struct {
	value interface{}
	err   error
}

// memoRequest is a message to a memo server, telling what key client wants to get
// and what callback must be executed in order to get the result (if there is yet
// nothing in cache.
type memoRequest struct {
	key      string
	callback MemoCallback
	response chan<- memoResult
}

// memo is an implementation of a IMemo interface.
type memo struct {
	requests chan memoRequest
}

// NewMemo provides a concurrency-safe non-blocking memoization of a function.
// Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
func NewMemo() *memo {
	m := &memo{
		requests: make(chan memoRequest),
	}
	go m.server()

	return m
}

// Get takes the value for a key from cache, and if there is no such,
// executes the callback and puts the result into the cache.
func (m *memo) Get(key string, callback MemoCallback) (interface{}, error) {
	response := make(chan memoResult)
	m.requests <- memoRequest{
		key:      key,
		callback: callback,
		response: response,
	}
	res := <-response

	return res.value, res.err
}

// handles requests
func (m *memo) server() {
	mem := make(map[string]*memoEntry)
	for req := range m.requests {
		e := mem[req.key]
		if e == nil || e.res.err != nil {
			e = &memoEntry{
				ready: make(chan struct{}),
			}
			mem[req.key] = e

			go e.call(req.callback)
		}
		go e.deliver(req.response)
	}
}

// memoEntry is a struct to store callback execution result,
// notifies listeners that callback has been executed.
type memoEntry struct {
	res   memoResult
	ready chan struct{} // closed when res is ready
}

// calls a callback function and notifies listeners when result is ready
func (e *memoEntry) call(f MemoCallback) {
	e.res.value, e.res.err = f()
	close(e.ready)
}

// sends result to the response channel when result is ready
func (e *memoEntry) deliver(response chan<- memoResult) {
	<-e.ready
	response <- e.res
}
