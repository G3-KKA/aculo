package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BatchProviderTestSuite struct {
	suite.Suite
}

func TestBatchProvider(t *testing.T) {
	suite.Run(t, new(BatchProviderTestSuite))
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
	t.Len(provider.inUse, 10)
	t.Len(provider.batches, 10)
	t.Len(provider.batches[0], 100)
}
func (t *BatchProviderTestSuite) Test_allocBatches() {
	testSlice := make([][]domain.Event, 10)
	t.Empty(testSlice[0])
	allocBatches(testSlice, 122)
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
	t.True(provider.inUse[0].Load())
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
	t.False(provider.inUse[0].Load())
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
	for d := range 3000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//time.Sleep(time.Duration(d) * time.Millisecond * 10)
			batch, mustcall := provider.GetBatch()
			defer mustcall()
			t.NotNil(batch)
			t.Len(batch, 100)
			for i := range len(batch) {
				batch[i].Data = []byte{}
				batch[i].EID = "I am Goroutine: " + strconv.Itoa(d)
				batch[i].ProviderID = "__"
				batch[i].SchemaID = "__"
				batch[i].Type = "test"
			}
		}()
	}
	wg.Wait()
	//t.T().Logf("%#v", provider)
	// batch, returnToProviderFunc := provider.GetBatch()
	// t.True(provider.inUse[0].Load())
	// t.Len(batch, 100)
	// t.Len(provider.batches[0], 100)
	// event := domain.Event{
	// 	EID:        "1",
	// 	ProviderID: "1",
	// 	SchemaID:   "1",
	// 	Type:       "1",
	// 	Data:       []byte("1"),
	// }
	// batch[0] = event
	// t.Equal(batch[0], provider.batches[0][0])
	// returnToProviderFunc()
	// t.False(provider.inUse[0].Load())
}
