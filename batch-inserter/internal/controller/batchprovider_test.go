package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"context"
	"strconv"
	"sync"
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

	provider := NewBatchProvider[domain.Event](context.TODO(), config.Config{
		Broker: config.Broker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	wg := sync.WaitGroup{}
	go func() {
		for {
			time.Sleep(300 * time.Millisecond) // global run q !!!!!!
			t.T().Log(len(provider.batches))
		}

	}()
	for d := range 4000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(d) * time.Millisecond * 10)
			for range 100 {
				//time.Sleep(time.Duration(d) * time.Nanosecond)
				func() {
					batch, mustcall := provider.GetBatch()
					t.NotNil(batch)
					t.Len(batch, 100)
					for i := range len(batch) {
						batch[i].Data = []byte{}
						batch[i].EID = "I am Goroutine: " + strconv.Itoa(d)
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
	// t.T().Log(len(provider.batches))
}

func (t *BatchProviderTestSuite) Test_NewBatchProvider_AllocationLogic() {
	provider := NewBatchProvider[domain.Event](context.TODO(), config.Config{
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
	testSlice := make([][]domain.Event, 10)
	t.Empty(testSlice[0])
	allocUnderlyingBatch(testSlice, 122)
	t.Len(testSlice[0], 122)
	t.Len(testSlice, 10)
	t.Equal(domain.Event{}, testSlice[0][0])
}

func (t *BatchProviderTestSuite) Test_GetBatch() {
	provider := NewBatchProvider[domain.Event](context.TODO(), config.Config{
		Broker: config.Broker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	batch, returnToProviderFunc := provider.GetBatch()
	t.True(provider.used[0].Load())
	t.Len(batch, 100)
	t.Len(provider.batches[0], 100)
	event := domain.Event{
		EID:        "1",
		ProviderID: "1",
		SchemaID:   "1",
		Type:       "1",
		Data:       []byte("1"),
	}
	batch[0] = event
	t.Equal(batch[0], provider.batches[0][0])
	returnToProviderFunc()
	t.False(provider.used[0].Load())
}
