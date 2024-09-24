package logger

import "go.uber.org/zap"

// Attaches specific name to zap.AtomicLevel
type NamedLevel struct {
	zap.AtomicLevel
	name string
}

// Returns name of
func (lvl NamedLevel) Name() string {
	return lvl.name
}

func withName(name string, level zap.AtomicLevel) NamedLevel {
	return NamedLevel{level, name}
}
