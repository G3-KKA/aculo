package logger

import (
	"io"
	"master-service/internal/config"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ io.Writer = &Logger{}

type (
	// Logger as object, no need for interfaces
	Logger struct {
		*zap.SugaredLogger
		levels []NamedLevel
	}
)

// # Create new logger
func New(config config.Logger) (Logger, error) {

	levels := make([]NamedLevel, 0, len(config.Cores))
	cores := make([]zapcore.Core, 0, len(config.Cores))

	// Iterating thorough config cores and creating zapcore.Cores out of them
	for _, core := range config.Cores {
		logDest, err := assembleDestination(string(core.Path))
		if err != nil {
			if core.MustCreateCore {
				return Logger{}, err
			}
			continue
		}
		encoder, err := setEncoder(core.EncoderLevel)
		if err != nil {
			return Logger{}, err
		}
		namedLevel := withName(core.Name, zap.NewAtomicLevelAt(zapcore.Level(core.Level)))
		levels = append(levels, namedLevel)
		cores = append(cores, zapcore.NewCore(
			encoder,               // production or development
			logDest,               // file or stderr/stdout // TODO Add remote dest support
			levels[len(levels)-1], // last level, every time
		))
	}
	levels = slices.Clip(levels)
	cores = slices.Clip(cores)
	if len(cores) == 0 {
		return Logger{}, ErrNoCoresWasInitialized
	}

	// Creating Sugar Logger from cores
	core := zapcore.NewTee(cores...)
	zaplogger := zap.New(core).Sugar()

	// First log message
	// That tells us that logger construction succeeded
	logger := Logger{
		SugaredLogger: zaplogger,
		levels:        levels,
	}
	logger.Debug("Logger construction succeeded")

	// TODO utilise returning stopFunc
	_ = syncOnTimout(logger, config.SyncTimeout)

	return logger, nil
}

// Writes to debug!
func (l *Logger) Write(in []byte) (n int, err error) {
	// l.Debug(string(in[:len(in)]))
	l.Debug(string(in[:]))
	return len(in), nil
}

// Return noop logger
func Noop() Logger {
	return Logger{
		SugaredLogger: zap.NewNop().Sugar(),
		levels:        []NamedLevel{},
	}
}

// ===================================================================================================================
