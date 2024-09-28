package logger

import "errors"

var (
	ErrNoCoresWasInitialized = errors.New("no cores was initialized")
	ErrCantOpenLogfile       = errors.New("unsuccessful logger core initialization, cant open log file")
	ErrUnknownEncoder        = errors.New("unknown encoder")
)
