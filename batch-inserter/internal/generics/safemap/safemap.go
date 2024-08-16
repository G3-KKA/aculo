package safemap

import (
	"errors"
	"sync"
)

func Make[K comparable, V any](cap int) *Map[K, V] {
	return &Map[K, V]{mx: &sync.RWMutex{}, mmap: make(map[K]V, cap)}
}

// Thread-safe generic-powered map
type Map[K comparable, V any] struct {
	mmap map[K]V

	mx *sync.RWMutex
}

// Thread safe adding to map
func (c *Map[K, V]) Store(key K, value V) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.mmap[key] = value
}
func (c *Map[K, V]) Load(key K) (V, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	value, ok := c.mmap[key]
	if !ok {
		return value, ErrValueNotFound
	}
	return value, nil

}

var ErrValueNotFound = errors.New("value not found")
