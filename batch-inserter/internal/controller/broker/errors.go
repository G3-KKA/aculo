package broker

import "errors"

var (
	ErrAlreadyHandling = errors.New("provided topic already handling")
	ErrStillHandling   = errors.New("provided topic still handled by not shutted down handler")
	ErrNotHandling     = errors.New("provided topic not handling")
)
