package unifaces

import (
	"errors"
)

// Unified interface-container for operation(s) happening previous to any specific api calls
//
// Treat it like http.Middleware or grpc.Interceptor
type Do[API any] interface {
	// Operation(s) happening previous to any specific api calls
	Api() (API, error)
}

// Unified interface-container for operation(s) happening previous to any specific api calls
//
// # With support for explicit client-size transaction closing
//
// Treat it like http.Middleware or grpc.Interceptor
type Tx[API any] interface {
	// Operation(s) happening previous to any specific api calls
	Tx() (API, TxClose, error)
}

// Explicit returned (errorproof) callback
type TxClose func() error

// Logic which should happen previous to any intercation with T
type InsideTxFunc func() error

// Generic Wrapper
type TxWrapper[T any] struct {
	Inside  []InsideTxFunc
	TxClose TxClose
	T       T
}

// # Generic wrapper, its enough for simple logic, like counter++ or atomic.CAS()
//
// If transaction internals requires complex logic
// , its better to implement your own T.Tx()
func WithTx[T any](t T, txclose TxClose, inside ...InsideTxFunc) Tx[T] {
	wrapper := TxWrapper[T]{
		Inside:  inside,
		TxClose: txclose,
		T:       t,
	}
	return &wrapper
}

func (wrapper *TxWrapper[T]) Tx() (T, TxClose, error) {
	for _, f := range wrapper.Inside {
		err := f()
		if err != nil {
			var genericZeroValue T
			err = errors.Join(ErrTxInternalError, err)
			return genericZeroValue, func() error { return err }, err
		}
	}
	return wrapper.T, wrapper.TxClose, nil
}

var (
	ErrTxAlreadyClosed = errors.New("this transaction already closed")
	ErrTxInternalError = errors.New("error happened inside Tx() body")
)
