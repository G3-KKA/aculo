package controller

import (
	"aculo/batch-inserter/internal/config"
	"context"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type cacheLineOptimizedAtomicBool struct {
	atomic.Bool
	_ [64 - unsafe.Sizeof(atomic.Bool{})]byte
}

// Generic thread safe storage for Rewriteable Slices
type BatchProvider[T any] struct {
	inUse   []cacheLineOptimizedAtomicBool
	batches [][]T

	// preallocSize int
	batchSize int

	reallocating atomic.Bool

	mx *sync.RWMutex
}

func NewBatchProvider[T any](ctx context.Context, config config.Config) *BatchProvider[T] {
	bp := &BatchProvider[T]{
		// Evety atomic bool corresponds to one preallocated batch
		inUse:   make([]cacheLineOptimizedAtomicBool, config.BatchProvider.PreallocSize),
		batches: make([][]T, config.BatchProvider.PreallocSize),

		// preallocSize: config.BatchProvider.PreallocSize,
		batchSize: config.Broker.BatchSize,
		mx:        &sync.RWMutex{},
	}
	allocBatches(bp.batches, config.Broker.BatchSize)

	return bp

}

// # Using batch after [MustReturnFunc] will cause guaranteed datarace
//
// # Safe to call multiple times
//
// Not calling it will leak batches
type MustReturnFunc func()

// Thread safe, allocation free
// Every batch is preallocated til its , do not try to append to it it will cause immediate reallocation
func (p *BatchProvider[T]) GetBatch() ([]T, MustReturnFunc) {
	var returnCalled atomic.Bool
	p.mx.RLock()
	if len(p.batches) > 200 {
		func() {

		}()
	}
	for i := range len(p.inUse) {

		// Found a free batch
		if p.inUse[i].CompareAndSwap(false, true) {
			defer p.mx.RUnlock()
			f := func() {

				// Multi-Call safe measure
				if returnCalled.CompareAndSwap(false, true) {
					p.inUse[i].Store(false)
				}
			}
			return p.batches[i], f
		}
	}
	p.mx.RUnlock()
	// If someone alredy reallocating, wait for it
	// TODO Need to test for dataraces, bu seems legit for me
	if !p.reallocating.CompareAndSwap(false, true) {
		time.Sleep(10 * time.Millisecond)
		return p.GetBatch()
	}
	p.mx.Lock()
	prevSize := len(p.batches[0])
	multiplier := 2
	p.inUse = make([]cacheLineOptimizedAtomicBool, len(p.inUse)*multiplier)
	p.batches = make([][]T, len(p.batches)*multiplier)

	allocBatches(p.batches, prevSize)
	//p.inUse[len(p.inUse)-1].Store(true)
	/* 	f := func() {

	if returnCalled.CompareAndSwap(false, true) {
		p.inUse[len(p.inUse)-1].Store(false)
		}
		}
	*/
	p.reallocating.Store(false)
	p.mx.Unlock()
	return p.GetBatch()
}
func allocBatches[T any](batches [][]T, size int) {
	for i := range len(batches) {
		batches[i] = make([]T, size)
	}
}
