package txface

import (
	"errors"
	"sync"
)

// Logic which should happen previous to any intercation with T
type InsideFunc func() error

// Generic Wrapper
type TxWrapper[T any] struct {
	Inside []InsideFunc
	Commit Commit
	T      T
}

// # Generic wrapper, its enough for simple logic, like counter++ or atomic.CAS()
//
// If transaction internals requires complex logic
// , its better to implement your own T.Tx()
func WithTx[T any](t T, commit Commit, inside ...InsideFunc) Tx[T] {
	wrapper := TxWrapper[T]{
		Inside: inside,
		Commit: commit,
		T:      t,
	}
	return &wrapper
}

func (wrapper *TxWrapper[T]) Tx() (T, Commit, error) {
	for _, f := range wrapper.Inside {
		err := f()
		if err != nil {
			var genericZeroValue T
			err = errors.Join(ErrTxInternalError, err)
			return genericZeroValue, func() error { return err }, err
		}
	}
	return wrapper.T, wrapper.Commit, nil
}
func WithLocker(lock sync.Locker) InsideFunc {
	f := func() error {
		lock.Lock()
		return nil
	}
	return InsideFunc(f)
}

// Example
func _() {
	var a = struct {
		val  int
		flag bool
	}{
		val:  0,
		flag: false,
	}
	var mx sync.Mutex
	commit := Commit(func() error {
		mx.Unlock()
		return nil
	})
	txwrapped := WithTx(&a, commit, WithLocker(&mx))

	//
	//
	a2, commit2, err := txwrapped.Tx()
	if err != nil {
		return
	}
	defer commit2()
	a2.flag = true

}
