package controller

import (
	"aculo/batch-inserter/internal/config"
	"context"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Generic, thread safe storage for slices to rewrite
type BatchProvider[T any] struct {

	// Slice of atomic bool's which index corresponds to one batch
	// len(used) == len(batches)
	used []atomic.Bool

	//
	batches [][]T

	reallocating atomic.Bool

	// RLock() in this case means that we are just trying to get batch
	// Lock()  in the other hand means goroutine that will reallocate [used] and [batches]
	//         while realloc happens this G needs serialised RW access to [used] and [batches]
	mx *sync.RWMutex
}

func NewBatchProvider[T any](ctx context.Context, config config.Config) *BatchProvider[T] {
	bp := &BatchProvider[T]{
		// Evety atomic bool corresponds to one preallocated batch
		used:    make([]atomic.Bool, config.BatchProvider.PreallocSize),
		batches: make([][]T, config.BatchProvider.PreallocSize),

		mx: &sync.RWMutex{},
	}
	allocUnderlyingBatch(bp.batches, config.Broker.BatchSize)

	return bp

}

// # Using batch after [MustReturnFunc] will cause guaranteed datarace
//
// # Safe to call multiple times
//
// Not calling it, batch will be blocked indefinitely
type MustReturnFunc func()

// # Get batch to rewrite
//
// Thread safe
func (p *BatchProvider[T]) GetBatch() ([]T, MustReturnFunc) {
	// This atomic used to ensure that used falg will be set by client to [false] exaclty once
	var returnCalled atomic.Bool

	// Potentially gopark() if reallocating
	p.mx.RLock()

	// Search through every atomic
	for i := range len(p.used) {

		// Try acquire
		if p.used[i].CompareAndSwap(false, true) {
			defer p.mx.RUnlock()
			f := func(used *atomic.Bool) func() {
				return func() {
					// Multiply f() call's safe measure, only first will store false
					if returnCalled.CompareAndSwap(false, true) {
						used.Store(false)
					}
				}
			}(&p.used[i])
			return p.batches[i], f
		}
	}
	// Wait for all G's in search cycle
	p.mx.RUnlock()
	// If someone alredy reallocating, wait for it
	// TODO Need to test for dataraces, bu seems legit for me
	if !p.reallocating.CompareAndSwap(false, true) {
		return p.GetBatch()
	}
	p.mx.Lock()
	prevSize := len(p.batches[0])
	multiplier := 2
	p.used = make([]atomic.Bool, len(p.used)*multiplier)
	p.batches = make([][]T, len(p.batches)*multiplier)

	allocUnderlyingBatch(p.batches, prevSize)
	go func() {
		os.Stderr.WriteString("imherer\n")
		time.Sleep(300 * time.Millisecond) //  ЗДЕСЬ НАС ВЫТЕСНЯЮТ В ГЛОБАЛ РАН Q ???
		p.reallocating.Store(false)
		p.mx.Unlock()
	}()
	//time.Sleep(time.Duration(300) * time.Millisecond)
	return p.GetBatch()
}

func allocUnderlyingBatch[T any](batches [][]T, size int) {
	for i := range len(batches) {
		batches[i] = make([]T, size)
	}
}
