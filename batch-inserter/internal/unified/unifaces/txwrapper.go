package unifaces

import "errors"

// Logic which should happen previous to any intercation with T
type InsideTxFunc func() error

// Generic Wrapper
type TxWrapper[T any] struct {
	TxInside []InsideTxFunc
	TxClose  TxClose
	T        T
}

// # Generic wrapper, its enough for simple logic, like counter++ or atomic.CAS()
//
// If transaction internals requires complex logic
// , its better to implement your own T.Tx()
func WithTx[T any](t T, txclose TxClose, txinside ...InsideTxFunc) Tx[T] {
	wrapper := TxWrapper[T]{
		TxInside: txinside,
		TxClose:  txclose,
		T:        t,
	}
	return &wrapper
}

func (wrapper *TxWrapper[T]) Tx() (T, TxClose, error) {
	for _, f := range wrapper.TxInside {
		err := f()
		if err != nil {
			var genericZeroValue T
			err = errors.Join(ErrTxInternalError, err)
			return genericZeroValue, func() error { return err }, err
		}
	}
	return wrapper.T, wrapper.TxClose, nil
}
