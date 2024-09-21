package asyncerror

import (
	"errors"
	"sync"
)

// # Sometimes code has to asynchronously modify( Join ) the same error.
//
// # This wrapper solves the problem via simple rwmutex
type AsyncError struct {
	err error
	mx  sync.RWMutex
}

// Get error under read lock,
// it might be invalid right after func returns
func (asyncErr *AsyncError) Err() error {
	asyncErr.mx.RLock()
	defer asyncErr.mx.RUnlock()
	return asyncErr.err
}

// # Safe to call asynchorously
func (asyncErr *AsyncError) Join(errs ...error) {
	asyncErr.mx.Lock()
	errs = append(errs, asyncErr.err)
	asyncErr.err = errors.Join(errs...)
	asyncErr.mx.Unlock()
}

// An mutex.RUnlock returned to the client
type ClientsideUnlock func()

// Get error under read lock,
// it is invalid only after client calls [ClientsideUnlock],
// returned by the function.
//
// # Not calling returned function might result in deadlock!
//
// # Its better to never use it, otherwise be careful
func (asyncErr *AsyncError) ErrClientsideUnlock() (ClientsideUnlock, error) {
	asyncErr.mx.RLock()
	return asyncErr.mx.RUnlock, asyncErr.err
}
