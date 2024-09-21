package option

type OptionFunc[T any] func(*T) error
