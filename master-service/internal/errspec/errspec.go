// # This package provides wrappers around errors.Join/New
//
// Idea that you first declare generic named error, that you CAN identify
// using errors.Is, but have no context, and for logs you provide this context.
package errspec

import (
	"errors"
	"fmt"
)

// Any context-less error that needs to be clarified.
type NotSpecifiedError error

// Same joins generic error with specific of 2 value being the same.
func Same(err NotSpecifiedError, a, b any) error {
	formatted := fmt.Sprintf(" |%v|  same as  |%v| ", a, b)

	return errors.Join(err, errors.New(formatted))
}

// Msg adds message to generic error.
func Msg(err NotSpecifiedError, msg string) error {
	return errors.Join(err, errors.New(msg))
}

// MsgValue adds message and value to generic error.
func MsgValue(err NotSpecifiedError, msg string, v any) error {
	formatted := fmt.Sprintf("%s:\t %v", msg, v)

	return errors.Join(err, errors.New(formatted))
}
