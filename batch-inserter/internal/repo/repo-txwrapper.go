package repository

import (
	"aculo/batch-inserter/internal/interfaces/shuttable"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/unified/unierrors"
	"sync"
	"sync/atomic"
)

// Tx wrapper over [repoapi]
type repository struct {
	r *repoapi

	unavailable atomic.Bool
	mx          sync.RWMutex
}

// # Common middleware for all api calls
//
// Safe to close multiple times,
//
// Safe to close if error happened in Tx itself, but not necessary.
//
// In case of error [txface.Commit] does nothing but return the same error.
func (rtx *repository) Tx() (*repoapi, txface.Commit, error) {

	// Immediate return if api was already shutted down
	if rtx.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	//
	// Acquire transaction
	rtx.mx.RLock()

	// # Rare case
	//
	// If someone sets unavaliable = true AND acquire  Lock()
	// in moment of time when we already checked unavaliable == false
	// and still not acquired RLock() we could, after Unlock() called
	// -> access memory that already "shutted down"
	if rtx.unavailable.Load() {
		rtx.mx.RUnlock()
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	commit := sync.OnceValue(func() error {
		rtx.mx.RUnlock()
		return nil
	})
	return rtx.r, txface.Commit(commit), nil

}

// [shuttable.Shuttable]
func (rtx *repository) ShuttedDown() bool {
	rtx.mx.Lock()
	defer rtx.mx.Unlock()
	return rtx.unavailable.Load()
}

// [shuttable.Shuttable]
func (rtx *repository) Shutdown() (err error) {
	rtx.mx.Lock()

	if !rtx.unavailable.CompareAndSwap(false, true) {
		return shuttable.ErrAlreadyShuttnigDown
	}

	defer rtx.mx.Unlock()
	err = rtx.r.shutdown()
	return
}

// [txface.Tx].[*brokerapi] --> [txface.Tx].[BrokerAPI]
func (rtx *repository) WrapAPI() txface.Tx[RepositoryAPI] {
	return &wrap{
		rtx: rtx,
	}

}

// [txface.ApiWrapper].[BrokerAPI]
type wrap struct {
	rtx *repository
}

// [txface.Tx].[BrokerAPI]
func (w *wrap) Tx() (RepositoryAPI, txface.Commit, error) {
	return w.rtx.Tx()
}
