package logger

import "go.uber.org/zap"

// zap.AtomicLevel with specific name.
type NamedLevel struct {
	zap.AtomicLevel
	name string
}

func withName(name string, level zap.AtomicLevel) NamedLevel {
	return NamedLevel{level, name}
}

// Name returns name of the logging level.
func (lvl NamedLevel) Name() string {
	return lvl.name
}
