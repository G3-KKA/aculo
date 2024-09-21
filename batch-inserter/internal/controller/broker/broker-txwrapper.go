package broker

import (
	"aculo/batch-inserter/internal/interfaces/shuttable"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/unified/unierrors"
	"sync"
	"sync/atomic"
)

// Tx wrapper over [brokerapi]
type broker struct {
	b *brokerapi

	unavailable atomic.Bool
	mx          sync.RWMutex
}

// # Common middleware for all api calls
//
// Safe to close multiple times,
//
// Safe to close if error happened in tx itself, but not necessary.
//
// In case of error [txface.Commit] does nothing but return the same error.
func (btx *broker) Tx() (*brokerapi, txface.Commit, error) {

	// Immediate return if api was already shutted down
	if btx.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	//
	// Acquire transaction
	btx.mx.RLock()

	// # Rare case
	//
	// If someone sets unavaliable = true AND acquire  Lock()
	// in moment of time when we already checked unavaliable == false
	// and still not acquired RLock() we could, after Unlock() called
	// -> access memory that already "shutted down"
	if btx.unavailable.Load() {
		btx.mx.RUnlock()
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	commit := sync.OnceValue(func() error {
		btx.mx.RUnlock()
		return nil
	})
	return btx.b, txface.Commit(commit), nil

}

// [shuttable.Shuttable]
func (btx *broker) ShuttedDown() bool {
	btx.mx.Lock()
	defer btx.mx.Unlock()
	return btx.unavailable.Load()
}

// [shuttable.Shuttable]
func (btx *broker) Shutdown() (err error) {
	btx.mx.Lock()

	if !btx.unavailable.CompareAndSwap(false, true) {
		return shuttable.ErrAlreadyShuttnigDown
	}

	defer btx.mx.Unlock()
	err = btx.b.shutdown()
	return
}

// [txface.Tx].[*brokerapi] --> [txface.Tx].[BrokerAPI]
func (btx *broker) WrapAPI() txface.Tx[BrokerAPI] {
	return &wrap{
		btx: btx,
	}

}

// [txface.ApiWrapper].[BrokerAPI]
type wrap struct {
	btx *broker
}

// [txface.Tx].[BrokerAPI]
func (w *wrap) Tx() (BrokerAPI, txface.Commit, error) {
	return w.btx.Tx()
}
