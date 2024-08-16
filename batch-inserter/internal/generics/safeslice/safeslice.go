package safeslice

import "sync"

type Slice[T any] struct {
	data []T
	mx   *sync.RWMutex
}

func Make[T any](len uint, cap uint) *Slice[T] {
	return &Slice[T]{

		data: make([]T, len, cap),
		mx:   &sync.RWMutex{},
	}
}
func (ss *Slice[T]) Append(args ...T) {

	ss.mx.Lock()
	defer ss.mx.Unlock()

	ss.data = append(ss.data, args...)

}
func (ss *Slice[T]) GetRead() (data []T, readUnlock func()) {

	ss.mx.RLock()

	return ss.data, ss.mx.RUnlock

}
func (ss *Slice[T]) Len(*Slice[T]) (length int, readUnlock func()) {

	ss.mx.RLock()

	return len(ss.data), ss.mx.RUnlock
}
func (ss *Slice[T]) Cap() (capacity int, readUnlock func()) {

	ss.mx.RLock()

	return cap(ss.data), ss.mx.RUnlock

}
