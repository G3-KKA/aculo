package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"context"
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
