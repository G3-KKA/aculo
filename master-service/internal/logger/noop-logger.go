package logger

import "log"

var _ ILogger = &nooplogger{}

type nooplogger struct {
}

// Debug implements Logger.
func (n *nooplogger) Debug(args ...any) {}

// Error implements Logger.
func (n *nooplogger) Error(args ...any) {}

// Fatal implements Logger.
func (n *nooplogger) Fatal(args ...any) {
	log.Fatal(args...)
}

// Info implements Logger.
func (n *nooplogger) Info(args ...any) {}

func NewNoopLogger() *nooplogger {
	return &nooplogger{}
}
