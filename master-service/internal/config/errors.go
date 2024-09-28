package config

import "errors"

var (
	ErrZeroValueConfig = errors.New("zero value config")
	ErrEnvNotDefined   = errors.New("env not defined")
)
