package controller

import (
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/controller"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

const bench_topic = "bench_topic"

var address []string = []string{"localhost:9092", "localhost:9093"}

var reqdata = []byte(`{"id":1, "name":"test"}`)

func BenchmarkSample(b *testing.B) {
	viper.BindEnv("BENCH_CONTROLLER")
	if viper.GetString("BENCH_CONTROLLER") == "" {
		b.Skip("skipping controller benchmark, env not set")
	}

	mock_logger := logger.NewMockLogger(b)
	mock_logger.On("Debug", mock.Anything).Run(func(args mock.Arguments) {
		b.Log(args)
	})
	mock_logger.On("Debug", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		b.Log(args)
	})
	// mock_logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
	// 	t.T().Log(args)
	// })
	// mock_logger.On("Info", mock.Anything).Run(func(args mock.Arguments) {
	// 	t.T().Log(args)
	// })
	//mock_logger.On("Info", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
	//	b.Log(args)
	//})
	//mock_logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
	//	b.Log(args)
	//})

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
			Topic:     bench_topic,
			BatchProvider: config.BatchProvider{
				PreallocSize: 10,
			},
		},
		Repository: config.Repository{},
	}
	ctl, err := controller.New(context.Background(), cfg, mock_logger, mock_service)
	for i := range 12 {

		f := func(ib int) {
			counter := 0
			go func() {
				for {
					time.Sleep(time.Second * 5)
					b.Logf("me:%d, already send: %d", ib, counter)
				}
			}()

			producer, err := sarama.NewAsyncProducer(cfg.Broker.Addresses, nil)

			if err != nil {
				b.Error(err.Error())
			}
			defer producer.Close()
			msg := sarama.ProducerMessage{
				Topic: bench_topic,
				Value: sarama.ByteEncoder(reqdata),
			}
			for {
				counter++

				if counter%200 == 0 {
					b.Logf("me:%d, already send: %d", ib, counter)
				}
				producer.Input() <- &msg
				if err != nil {
					b.Error(err.Error())
				}
			}
		}
		go f(i)
	}
	if err != nil {
		b.Error(err.Error())
	}

	for i := 0; i < b.N; i++ {
		f := func() {
			ctlapi, txclose, err := ctl.Tx()
			if err != nil {
				b.Error(err.Error())
			}
			defer txclose()
			// КАК ЭТО ВЫКЛЮЧАТЬ И ИЛИ БЕНЧМАРКАТЬ !!!!!!!*********************
			err = ctlapi.HandleBatch(context.Background())
			if err != nil && err != unierrors.ErrOperationInterrupted {
				b.Error(err.Error())
			}
		}
		b.Log("hello motherfuckers")
		f()

	}
}
