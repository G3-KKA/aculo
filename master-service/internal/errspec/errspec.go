// # This package provides wrappers around errors.Join/New
//
// Idea that you first declare generic named error, that you CAN identify
// using errors.Is, but have no context, and for logs you provide this context
package errspec

import (
	"errors"
	"fmt"
)

type NotSpecifiedError error

// Joins generic error, that specify what happened

func Same(err NotSpecifiedError, a, b any) error {
	samestring := fmt.Sprintf(" |%v|  same as  |%v| ", a, b)
	err = errors.Join(err, errors.New(samestring))
	return err
}
