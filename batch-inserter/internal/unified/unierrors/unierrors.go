package unierrors

import "errors"

var (
	ErrUnavailable          = errors.New("unavailable")
	ErrOperationInterrupted = errors.New("operation interrupted")
	//ErrAlreadyShuttingDown  = errors.New("already shutting down")
)
