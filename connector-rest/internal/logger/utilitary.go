package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: There might be problems with /stderr and debugging go code via deluge
func logFileFromPath(path string) (*os.File, error) {

	// Trying to create log file
	logfile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		// There is common case that directory doesn't exist
		// So we try to create it
		err = os.Mkdir(filepath.Dir(path), 0777)
		if err != nil {
			return nil, err
		}

		// Retry to create log file
		logfile, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("unsuccessful logger core initialization, cannot open log file")
		}
	}
	return logfile, nil
}

// Be careful when changing config.logger.cores.encoderLevel in runtime
// Might Panic!
func setEncoder(name string) (zapcore.Encoder, error) {

	if name == "production" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), nil
	}
	if name == "development" {
		return zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), nil
	}
	return nil, fmt.Errorf("unknown encoder: %s", name)
}

// Calls .Sync() for provided logger every syncTimeout
func syncOnTimout(logger *zap.SugaredLogger, syncTimeout time.Duration) (stop StopFunc) {
	go func() {
		ticker, stop1 := NewStopableTicker(syncTimeout)
		stop = stop1
		C := ticker.Chan()
		for {
			if ticker.Closed() {
				return
			}
			<-C
			_ = logger.Sync()
		}
	}()
	return stop
}

// G3KKA Template Library.utils
// Use it as you wish, even if i wrote tests -- you should never use it in production code

// Name speaks for itself
type StopableTicker struct {
	c      chan time.Time
	closed atomic.Bool
}

type StopFunc func()

// Incapsulates channel, provided channel is read-only and cannot be closed by hand
// It WILL be closed inside ticker logic, after calling StopFunc
func (ticker *StopableTicker) Chan() <-chan time.Time {
	return ticker.c
}

// Returns state of the ticker
func (ticker *StopableTicker) Closed() bool {
	return ticker.closed.Load()
}

// Uses time.Ticker inside, adds stopper functionality
// StopFunc closes the channel and stores true in closed flag
// False ticks may occur!
func NewStopableTicker(d time.Duration) (*StopableTicker, StopFunc) {
	stopable := &StopableTicker{
		c:      make(chan time.Time, 1),
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
