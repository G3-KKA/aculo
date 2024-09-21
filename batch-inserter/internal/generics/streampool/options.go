package streampool

import "aculo/batch-inserter/internal/interfaces/option"

type (
	Options struct {
		StartSize uint64
	}
	PoolOption option.OptionFunc[Options]
)

func DefaultOptions() Options {
	return Options{
		StartSize: 20,
	}
}
