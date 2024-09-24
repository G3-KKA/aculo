package logger

import (
	"master-service/internal/config"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate mockery --filename=mock_logger.go --name=Logger --dir=. --structname MockLogger  --inpackage=true
type (
	ILogger interface {
		Debug(args ...any)
		Info(args ...any)
		Error(args ...any)
		Fatal(args ...any)
	}
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

// ===================================================================================================================
