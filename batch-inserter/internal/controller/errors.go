package controller

import (
	"errors"
)

var (
	ErrFailedToCloseConsumer = errors.New("failed to close consumer")
)
