package streampool

import (
	"sync"
	"sync/atomic"
)

type (
	Pool struct {
		stop chan struct{}

		workers   map[string]*worker
		workersmx sync.Mutex // serialise access to workers map

		unavailable atomic.Bool
		wg          sync.WaitGroup
	}

	// Stop channels are unique for every worker
	//
	// That they may be closed separately
	PoolFunc func(stop <-chan struct{})
)

// This pool intended to be used with any task
// That WON'T close eventually and require trigger from the outside
func NewStreamPool(opts ...PoolOption) *Pool {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	return &Pool{
		stop:        make(chan struct{}),
		workers:     make(map[string]*worker, options.StartSize),
		workersmx:   sync.Mutex{},
		unavailable: atomic.Bool{},
		//mx:          sync.RWMutex{},
		wg: sync.WaitGroup{},
	}
}

// # Starts concurrently execute the function
//
// If Go returned with non'nil error it guaranties:
//
// - function was not called
//
// - state of the pool not changed
func (pool *Pool) Go(name string, f PoolFunc) (err error) {

	w := worker{
		name:               name,
		stop:               make(chan struct{}, cap(pool.stop)),
		closestop:          sync.Once{},
		triggerShuttedDown: make(chan struct{}),
	}
	err = pool.addWorker(&w, name)
	if err != nil {
		return
	}

	//
	poolShutdownTrigger := func() {
		select {
		case <-pool.stop: // Pool shutted down
			w.asyncStop()
		case <-w.triggerShuttedDown: // Worker shutted down first, prevent goroutine leak
		}

	}
	go poolShutdownTrigger()

	//
	//
	f2 := func() {

		// last step , signalise that everything

		defer pool.doneWorker()
		defer pool.deleteWorker(w.name)
		// shutdownwait need that
		w.do(f)
		// when worker returns from the funtction
		// -1- firstly it signalises about that via closing the trigger
		// -2- then it is deleted from map
		// -3- lastly counter in mutex decreased by one, which tells pool that one worker is done

	}
	go f2()

	return nil
}

// # Function returns immediately
//
// # If worker listening to the stop channel it will eventually stopped
//
// If error happened in worker shutting down -- client will not know
func (pool *Pool) AsyncStopWorker(name string) {
	go func() {
		_ = pool.StopWorkerWait(name)
	}()
}

// If worker listening to the stop channel it will eventually stopped
//
// # Function returns only after
func (pool *Pool) StopWorkerWait(name string) error {

	//
	// Fast return if pool already shutted down
	if pool.unavailable.Load() {
		return ErrPoolShuttedDhown
	}

	// Serialise access to workers map
	pool.workersmx.Lock()

	//
	// Slow return if function was called in moment of pool shutdown
	if pool.unavailable.Load() {
		pool.workersmx.Unlock()
		return ErrPoolShuttedDhown
	}

	worker, ok := pool.workers[name]
	if !ok {
		pool.workersmx.Unlock()
		return ErrWorkerNotFound
	}

	pool.workersmx.Unlock()

	worker.stopWait()

	return nil
}

// Every worker will be notified
//
// # Function exits only when every last worker returned from the task
func (pool *Pool) ShutdownWait() {
	if !pool.unavailable.CompareAndSwap(false, true) {
		return
	}

	close(pool.stop)

	pool.workersmx.Lock() // Берем серелизованный доступ к мапе
	for _, b := range pool.workers {
		b.asyncStop()
	}
	pool.workersmx.Unlock()

	pool.waitForWorkers()

}

func (pool *Pool) AsyncShutdown() {
	go func() {
		pool.ShutdownWait()

	}()
}

//
//
// Inernal functions
//
//

func (pool *Pool) waitForWorkers() {
	pool.wg.Wait()
}
func (pool *Pool) addWorker(w *worker, name string) error {

	// Add one to the worker counter
	pool.wg.Add(1)

	//
	// Critical section
	pool.workersmx.Lock()

	if pool.unavailable.Load() {
		pool.wg.Done()
		return ErrPoolShuttedDhown

	}
	if _, exist := pool.workers[name]; exist {
		return ErrWorkerAlreadyExist
	}
	pool.workers[name] = w

	pool.workersmx.Unlock()
	return nil
}
func (pool *Pool) doneWorker() {
	pool.wg.Done()
}
func (pool *Pool) deleteWorker(name string) {

	pool.workersmx.Lock()
	// delete is a no-op, so we can ommit check that worker still exists
	delete(pool.workers, name)
	pool.workersmx.Unlock()

}
