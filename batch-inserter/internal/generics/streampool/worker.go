package streampool

import "sync"

type worker struct {

	//
	// identifier of worker, should be unique !
	name string

	//
	// pool-side stop, unique per worker
	stop chan struct{}
	// used to close stop channel
	closestop sync.Once

	//
	// If closed --> worker returned from the task and deleted from map
	triggerShuttedDown chan struct{}
}

func (w *worker) asyncStop() {
	f := func() {
		close(w.stop)
	}
	w.closestop.Do(f)
}
func (w *worker) stopWait() {
	f := func() {
		close(w.stop)
	}
	w.closestop.Do(f)
	<-w.triggerShuttedDown

}
func (w *worker) do(f PoolFunc) {

	// defer in case of panic
	defer close(w.triggerShuttedDown)
	f(w.stop)
}
