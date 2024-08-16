package unifaces

import "errors"

// Unified interface-container for operation(s) happening previous to any specific api calls
//
// Treat it like http.Middleware or grpc.Interceptor
type Tx[API any] interface {
	// Operation(s) happening previous to any specific api calls
	Tx() (API, TxClose, error)
}

// Explicit returned (errorproof) callback
type TxClose func() error

var (
	ErrTxAlreadyClosed = errors.New("this transaction already closed")
)
