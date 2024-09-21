package shuttable

import "errors"

// Anything that can be shutdown, unable to restart
//
//go:generate mockery --filename=mock_shuttable.go --name=Shuttable --dir=. --structname MockShuttable  --inpackage=true
type Shuttable interface {
	// Should return [ErrAlreadyShuttnigDown] on second call.
	Shutdown() error
	ShuttedDown() bool
}

var (
	ErrAlreadyShuttnigDown = errors.New("already shutting down")
)
