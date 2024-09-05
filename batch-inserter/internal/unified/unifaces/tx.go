package unifaces

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
