package streampool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const noexpot_WORKER_COUNT_TEST = 1000

func TestPool(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 10)
	pool := NewStreamPool()
	wg := sync.WaitGroup{}
	for i := range noexpot_WORKER_COUNT_TEST {

		wg.Add(1)
		name := "testworker"
		counter := 0
		//	pool.AsyncStopWorker("masya")
		testfunc := func(stop <-chan struct{}) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				case <-ticker.C:
					t.Log("new tick")
					counter++
					if counter > 3 {
						time.Sleep(time.Millisecond * 100)
					}

				}
			}
		}
		err := pool.Go(name+fmt.Sprint(i), PoolFunc(testfunc))
		if err != nil {

			panic(err)
		}
	}
	time.Sleep(time.Millisecond * 500)
	for {
		pool.workersmx.Lock()
		if len(pool.workers) == noexpot_WORKER_COUNT_TEST {
			pool.workersmx.Unlock()
			err := pool.StopWorkerWait("testworker16")
			if err != nil {
				panic(err)
			}
			break
		}
		pool.workersmx.Unlock()
	}
	for {
		pool.workersmx.Lock()
		if len(pool.workers) == -1+noexpot_WORKER_COUNT_TEST {
			pool.workersmx.Unlock()
			err := pool.StopWorkerWait("testworker16")
			if err == nil {
				panic("err")
			}
			break
		}
		pool.workersmx.Unlock()
	}
	pool.ShutdownWait()
	wg.Wait()
	if len(pool.workers) != 0 {
		panic("early return ! ")
	}
	err := pool.Go("", func(stop <-chan struct{}) {})
	if err != ErrPoolShuttedDhown {
		panic("pool is actulally shutted down")
	}
	err = pool.StopWorkerWait("")
	if err != ErrPoolShuttedDhown {
		panic("pool is actulally shutted down")
	}

}
