package controller

import (
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/unified/unierrors"
	"context"
	"testing"
)

const bench_topic = "bench_topic"

var address []string = []string{"localhost:9092", "localhost:9093"}

func BenchmarkSample(b *testing.B) {

	mock_logger := logger.NewNoopLogger()

	mock_serviceapi := service.NewMockServiceAPI(b)

	mock_serviceapi.On("SendBatch", mock.Anything, mock.Anything).Return(nil)

	// mock_serviceapi.On("GracefulShutdown").Return(nil)

	// Mock Service
	mock_service := service.NewMockService(b)
	mock_service.On("Tx").Return(mock_serviceapi,
		unifaces.TxClose(func() error { return nil }),
		nil,
	)
	cfg := config.Config{
		Logger: config.Logger{},
		Broker: config.Broker{
			Addresses: address,
			BatchSize: 1000,
			Offset:    sarama.OffsetOldest, // we need to test only the throughput here
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
		Repository: config.Repository{},
	}
	ctl, err := controller.New(context.Background(), bench_topic, cfg, mock_logger, mock_service)
	if err != nil {
		b.Fatal(err.Error())

	}
	if err != nil {
		b.Error(err.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := func() {
			ctlapi, txclose, err := ctl.Tx()
			if err != nil {
				b.Error(err.Error())
			}
			defer txclose()
			err = ctlapi.HandleBatch(context.Background())
			if err != nil && err != unierrors.ErrOperationInterrupted {
				b.Error(err.Error())
			}
		}
		//b.Log("DEBUG_REMOVE_FOR_TRUE benchmark batch handled")
		f()

	}
}
