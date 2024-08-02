package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"context"
	"sync"
	"sync/atomic"
)

type BatchProvider struct {
	inUse   []atomic.Bool
	batches [][]domain.Event

	preallocSize int
	batchSize    int
	mx           *sync.RWMutex
}

func NewBatchProvider(ctx context.Context, config config.Config) *BatchProvider {
	bprov := &BatchProvider{
		inUse:   make([]atomic.Bool, config.BatchProvider.PreallocSize),
		batches: make([][]domain.Event, config.BatchProvider.PreallocSize),

		preallocSize: config.BatchProvider.PreallocSize,
		batchSize:    config.Brocker.BatchSize,
		mx:           &sync.RWMutex{},
	}
	allocBatches(bprov.batches, config.Brocker.BatchSize)

	return bprov

}
func allocBatches(batches [][]domain.Event, size int) {
	for i := range len(batches) {
		batches[i] = make([]domain.Event, size)
	}
}

type ReturnFunc func()

func (p *BatchProvider) GetBatch() ([]domain.Event, ReturnFunc) {

	p.mx.RLock()
	for i := range len(p.inUse) {
		if p.inUse[i].CompareAndSwap(false, true) {
			defer p.mx.RUnlock()
			return p.batches[i], func() { p.inUse[i].Store(false) }
		}
	}
	p.mx.RUnlock()

	p.mx.Lock()
	defer p.mx.Unlock()
	prevSize := len(p.batches[0])
	p.inUse = make([]atomic.Bool, len(p.inUse)*2)
	p.batches = make([][]domain.Event, len(p.batches)*2)

	allocBatches(p.batches, prevSize)

	p.inUse[len(p.inUse)-1].Store(true)
	return p.batches[len(p.batches)-1], func() { p.inUse[len(p.inUse)-1].Store(false) }
}
