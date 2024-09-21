package unierrors

import "errors"

var (
	ErrUnavailable          = errors.New("unavailable")
	ErrOperationInterrupted = errors.New("operation interrupted")
	ErrNotInitialisedYet    = errors.New("not initialised yet ")
	//ErrAlreadyShuttingDown  = errors.New("already shutting down")
)
