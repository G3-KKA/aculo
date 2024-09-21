package unierrors

import "errors"

var (
	ErrUnsuccessfulInitialisation = errors.New("unsuccessful logger initialisation")
	ErrLogWrittenPartially        = errors.New("log written partially")
)
