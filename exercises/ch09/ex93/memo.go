// Exercise 9.3: Extend the Func type and the (*Memo).Get method so that callers
// may provide an optional done channel through which they can cancel the
// operation (§8.9). The results of a cancelled Func call should not be cached.
//
// Package memo provides a concurrency-safe non-blocking memoization
// of a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a monitor goroutine.
package memo

import (
	"context"
	"fmt"
)

// Func is the type of the function to memoize.
type Func func(ctx context.Context, key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

// result wrapper, to store ready state
type entry struct {
	id    int
	res   result
	ctx   context.Context
	ready chan struct{}
}

// A request is a message requesting that the Func be applied to key.
type request struct {
	id       int
	key      string
	ctx      context.Context
	response chan<- result
}

type Memo struct {
	requests chan request
}

func New(f Func) *Memo {
	m := &Memo{requests: make(chan request)}
	go m.server(f)
	return m
}

func (memo *Memo) Get(ctx context.Context, key string, id int) (interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	response := make(chan result)
	memo.requests <- request{id, key, ctx, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	cacheEntry := func(e **entry, req request) {
		*e = &entry{id: req.id, ctx: req.ctx, ready: make(chan struct{})}
		cache[req.key] = *e
		go (*e).call(f, req.key)
	}
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key
			cacheEntry(&e, req)
		} else {
			<-e.ready
			if e.res.err != nil {
				select {
				case <-e.ctx.Done():
					// invalidate the previous cancelled cache
					fmt.Printf("[%d]:err %s (stale cache)\n", e.id, req.key)
					cacheEntry(&e, req)
				default:
				}
			}
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// Evaluate the function.
	e.res.value, e.res.err = f(e.ctx, key)
	// Broadcast the ready condition.
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}
