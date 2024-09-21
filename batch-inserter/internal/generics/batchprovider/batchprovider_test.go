package batchprovider

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/testing/asyncsuite"
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type BatchProviderTestSuite struct {
	suite.Suite
}

func TestBatchProvider(t *testing.T) {
	suite.Run(t, new(BatchProviderTestSuite))
}
func (t *BatchProviderTestSuite) Test_GetBatchAsync() {
	asyncT := asyncsuite.AsyncSuite(&t.Suite)

	provider := New[domain.Log](context.TODO(), config.Config{
		Broker: config.Broker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	mx := sync.Mutex{} // to eliminati datarace inside test's themselves !!!!
	done := atomic.Bool{}
	wgCleanup := sync.WaitGroup{}
	wgCleanup.Add(1)
	go func() {
		for !done.Load() {
			time.Sleep(300 * time.Millisecond)
			mx.Lock()
			asyncT.T().Log(len(provider.batches))
			mx.Unlock()
		}
		defer wgCleanup.Done()
	}()
	wg := sync.WaitGroup{}
	for d := range 4000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(d) * time.Millisecond)
			for range 100 {
				//time.Sleep(time.Duration(d) * time.Nanosecond)
				func() {
					mx.Lock()
					batch, mustcall := provider.GetBatch()
					mx.Unlock()
					asyncT.NotNil(batch)
					asyncT.Len(batch, 100)
					for i := range len(batch) {
						batch[i].Data = []byte{}
						batch[i].LogID = "I am Goroutine: " + strconv.Itoa(d)
						batch[i].ProviderID = "__"
						batch[i].SchemaID = "__"
						batch[i].Type = "test"
					}
					mustcall()
				}()
			}
		}()
	}
	wg.Wait()
	done.Store(true)
	wgCleanup.Wait()
	// t.T().Log(len(provider.batches))
}

func (t *BatchProviderTestSuite) Test_NewBatchProvider_AllocationLogic() {
	provider := New[domain.Log](context.TODO(), config.Config{
		Broker: config.Broker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	t.Len(provider.used, 10)
	t.Len(provider.batches, 10)
	t.Len(provider.batches[0], 100)
}
func (t *BatchProviderTestSuite) Test_allocBatches() {
	testSlice := make([][]domain.Log, 10)
	t.Empty(testSlice[0])
	allocUnderlyingBatch(testSlice, 122)
	t.Len(testSlice[0], 122)
	t.Len(testSlice, 10)
	t.Equal(domain.Log{}, testSlice[0][0])
}

func (t *BatchProviderTestSuite) Test_GetBatch() {
	provider := New[domain.Log](context.TODO(), config.Config{
		Broker: config.Broker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	batch, returnBatchFn := provider.GetBatch()
	t.True(provider.used[0].Load())
	t.Len(batch, 100)
	t.Len(provider.batches[0], 100)
	event := domain.Log{
		LogID:      "1",
		ProviderID: "1",
		SchemaID:   "1",
		Type:       "1",
		Data:       []byte("1"),
	}
	batch[0] = event
	t.Equal(batch[0], provider.batches[0][0])
	returnBatchFn()
	t.False(provider.used[0].Load())
}
