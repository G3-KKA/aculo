package httpctl

import "errors"

var (
	ErrShutdownTimeoutExceeded = errors.New("shutdown timeout exceeded")
)
