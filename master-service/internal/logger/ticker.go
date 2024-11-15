package logger

import (
	"sync/atomic"
	"time"
)

// Ticker without false-ticks.
type (
	stopableTicker struct {
		c      chan time.Time
		closed atomic.Bool
	}

	stopFunc func()
)

// Incapsulates channel, provided channel is read-only and cannot be closed by hand.
//
// It WILL be closed inside ticker logic, after calling StopFunc.
func (ticker *stopableTicker) Chan() <-chan time.Time {
	return ticker.c
}

// Returns state of the ticker.
func (ticker *stopableTicker) Closed() bool {
	return ticker.closed.Load()
}

// Uses time.Ticker inside, adds stopper functionality.
//
// StopFunc closes the channel and stores true in closed flag.
//
// False ticks may occur!
func newStopableTicker(d time.Duration) (*stopableTicker, stopFunc) {
	stopable := &stopableTicker{
		c:      make(chan time.Time),
		closed: atomic.Bool{},
	}
	go func() {
		ticker := time.NewTicker(d)
		for {
			val := <-ticker.C
			if stopable.closed.Load() {
				close(stopable.c)

				return
			}
			stopable.c <- val
		}
	}()
	stop := func() {
		stopable.closed.Store(true)
	}

	return stopable, stop
}
