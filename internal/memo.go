package internal

// MemoCallback is the type of the function to memoize.
type MemoCallback func() (interface{}, error)

// A result is the result of calling a MemoCallback.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// A request is a message requesting that the MemoCallback be applied to key.
type request struct {
	key      string
	callback MemoCallback
	response chan<- result // the client wants a single result
}

type Memo struct {
	requests chan request
}

// New returns a memoization of callback.  Clients must subsequently call Close.
func New(f MemoCallback) *Memo {
	memo := &Memo{
		requests: make(chan request),
	}
	go memo.server()
	return memo
}

func (memo *Memo) Get(key string, callback MemoCallback) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{
		key:      key,
		callback: callback,
		response: response,
	}
	res := <-response

	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server() {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{
				ready: make(chan struct{}),
			}
			cache[req.key] = e

			go e.call(req.callback) // call callback(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f MemoCallback) {
	e.res.value, e.res.err = f()
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}
