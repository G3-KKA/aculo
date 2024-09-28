package logger

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"master-service/internal/errspec"
)

const (
	FILE_PERMISSIONS      = 0o666
	DIRECTORY_PERMISSIONS = 0o777
)

// TODO: There might be problems with /stderr and debugging go code via dlv.
func assembleDestination(path string) (zapcore.WriteSyncer, error) {

	// Trying to create log file.
	logfile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, FILE_PERMISSIONS)
	if err != nil {
		// There is common case that directory doesn't exist,
		// So we try to create it.
		err = os.Mkdir(filepath.Dir(path), DIRECTORY_PERMISSIONS)
		if err != nil {
			return nil, err
		}

		// Retry to create log file.
		logfile, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, FILE_PERMISSIONS)
		if err != nil {
			return nil, errors.Join(ErrCantOpenLogfile, err)
		}
	}

	return logfile, nil
}

// Be careful when changing config.logger.cores.encoderLevel in runtime.
// Might Panic!
func setEncoder(name string) (zapcore.Encoder, error) {

	if name == "production" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), nil
	}
	if name == "development" {
		return zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), nil
	}

	return nil, errspec.MsgValue(ErrUnknownEncoder, "got encoder", name)
}

// Calls .Sync() for  every syncTimeout.
func syncOnTimout(logger Logger, syncTimeout time.Duration) (stop stopFunc) {
	go func() {
		var ticker *stopableTicker
		ticker, stop = newStopableTicker(syncTimeout)
		for {
			if ticker.Closed() {
				return
			}
			<-ticker.Chan()
			_ = logger.Sync()
		}
	}()

	return stop
}
