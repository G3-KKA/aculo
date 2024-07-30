package logger

import (
	"aculo/connector-restapi/internal/config"
	"fmt"
	"log"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Should be initialised via InitGlobalLogger()
var globalLogger *zap.SugaredLogger

// Little wrapper for future ease of identification
type LevelWithName struct {
	zap.AtomicLevel
	name string
}

// R_ONLY name
func (l LevelWithName) Name() string {
	return l.name
}

// []LevelWithName may be used to change specific output destination log levels
// Changing them in runtime is tread safe
func AssembleLogger(config config.Config) (*zap.SugaredLogger, []LevelWithName, error) {

	// May dynamicly change log levels in runtime, will be returned from InitLogger()
	levels := make([]LevelWithName, 0, len(config.Logger.Cores))

	// Creating cores fully dynamic from config
	// stderr/stdout supported, network not supported
	// TODO: Add network support
	cores := make([]zapcore.Core, 0, len(config.Logger.Cores))

	// Iterating thorough config cores and creating zapcore.Cores out of them
	for _, core := range config.Logger.Cores {
		logFile, err := logFileFromPath(string(core.Path))
		if err != nil {
			if core.MustCreateCore {
				return nil, nil, err
			}
			log.Println("One of the cores cannot be created, but it will be ignored")
			continue
		}
		encoder, err := setEncoder(core.EncoderLevel)
		if err != nil {
			return nil, nil, err
		}
		namedLevel := withName(core.Name, zap.NewAtomicLevelAt(zapcore.Level(core.Level)))
		levels = append(levels, namedLevel)
		cores = append(cores, zapcore.NewCore(
			encoder,               // production or development
			logFile,               // file or stderr/stdout
			levels[len(levels)-1], // last level, every time
		))
	}
	levels = slices.Clip(levels)
	cores = slices.Clip(cores)
	if len(cores) == 0 {
		return nil, nil, fmt.Errorf("logger initialization failed, has no cores")
	}

	// Creating zap.Cores
	// And merging them
	core := zapcore.NewTee(cores...)

	// Creating Logger from cores
	// And sugaring
	logger := zap.New(core)
	sugarlogger := logger.Sugar()

	// First log message
	// That tells us that logger construction succeeded
	defer sugarlogger.Sync()
	sugarlogger.Debug("Logger construction succeeded")

	return sugarlogger, levels, nil
}

// Useful for small apps where you want to log a bit
// Not sure about async
func InitGlobalLogger(config config.Config) (err error) {

	// Ignoring ability to change level in runtime for global usecase
	// TODO: it is not hard to add this feature, mb next time
	globalLogger, _, err = AssembleLogger(config)

	// Ignoring ablity to stop Sync'ing
	_ = syncOnTimout(globalLogger, config.Logger.SyncTimeout)

	return
}

// Please InitGlobalLogger first, thx
func Debug(args ...any) {
	globalLogger.Debug(args...)
}

// Please InitGlobalLogger first, thx
func Info(args ...any) {
	globalLogger.Info(args...)
}
func Fatal(args ...any) {
	globalLogger.Fatal(args...)
}
func Error(args ...any) {
	globalLogger.Error(args...)
}

func withName(name string, level zap.AtomicLevel) LevelWithName {
	return LevelWithName{level, name}
}
