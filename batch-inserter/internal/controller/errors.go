package controller

import (
	"errors"
)

var (
	ErrFailedToCloseConsumer          = errors.New("failed to close consumer")
	ErrRetryNotWorksConnectionRefused = errors.New("retry not works, connection refused")
)
