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
	provider := NewBatchProvider(context.TODO(), config.Config{
		Brocker: config.Brocker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	t.Equal(10, len(provider.inUse))
	t.Equal(10, len(provider.batches))
	t.Equal(100, len(provider.batches[0]))
}
func (t *BatchProviderTestSuite) Test_allocBatches() {
	testSlice := make([][]domain.Event, 10)
	t.Equal(0, len(testSlice[0]))
	allocBatches(testSlice, 122)
	t.Equal(10, len(testSlice))
	t.Equal(122, len(testSlice[0]))
	t.Equal(domain.Event{}, testSlice[0][0])
}

func (t *BatchProviderTestSuite) Test_GetBatch() {
	provider := NewBatchProvider(context.TODO(), config.Config{
		Brocker: config.Brocker{
			BatchSize: 100,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
	})
	batch, returnToProviderFunc := provider.GetBatch()
	t.Equal(true, provider.inUse[0].Load())
	t.Equal(100, len(batch))
	t.Equal(100, len(provider.batches[0]))
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
	t.Equal(false, provider.inUse[0].Load())
}
