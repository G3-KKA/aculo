package txface

import "errors"

var (
	ErrTxAlreadyClosed = errors.New("this transaction already closed")
	ErrTxInternalError = errors.New("error happened inside Tx() body")
)
