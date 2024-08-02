package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite_Intgr struct {
	suite.Suite
}

func TestController(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite_Intgr))
}

// ========================

func (t *ControllerTestSuite_Intgr) SetupSuite() {
	err := logger.InitGlobalLogger(config.Config{
		Logger: config.Logger{
			SyncTimeout: 300 * time.Millisecond,
			Cores: []config.LoggerCore{
				{
					Name:           "stdout",
					EncoderLevel:   "production",
					Path:           "/dev/stdout",
					Level:          0,
					MustCreateCore: true,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	// testutils // TODO testutils from connector

}
func (t *ControllerTestSuite_Intgr) BeforeTest(suiteName, testName string) {
	switch testName {
	default:
	}

}

func (t *ControllerTestSuite_Intgr) Test_Serve() {
	reqdata := []byte(`{"id":1, "name":"test"}`)
	go func() {
		producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
		if err != nil {
			panic(err)
		}
		defer producer.Close()
		msg := &sarama.ProducerMessage{
			Topic: "test_topic",
			Value: sarama.ByteEncoder(reqdata),
		}
		for range 3 {
			fmt.Println("send-0", string(reqdata))
			_, _, err = producer.SendMessage(msg)
			fmt.Println("send-1", string(reqdata))
			if err != nil {
				panic(err)
			}

		}

	}()
	mockservice := service.NewMockService(t.T())
	mockservice.On("SendBatch", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx == context.Background() || ctx == context.TODO()
	}), mock.MatchedBy(func(eventbatch []domain.Event) bool {
		for _, event := range eventbatch {
			if string(event.Data) != string(reqdata) {
				return false
			}
		}
		fmt.Println("success")
		return true
	})).Return(nil)
	//mock.AnythingOfType()
	//mock.MatchedBy()
	controller, err := New(context.Background(), config.Config{
		Brocker: config.Brocker{
			Addresses: []string{"localhost:9092"},
			Topic:     "test_topic",
			BatchSize: 3,
			BatchProvider: config.BatchProvider{
				PreallocSize: 2,
			},
		},
	}, mockservice)
	t.Nil(err)
	fmt.Println("controller service started", controller)
	err = controller.Serve()
	// ============================== TODO СДЕЛАТЬ ТАК ЧТОБЫ ЭТА ДУРА ВЫКЛЮЧАЛАСЬ ПО КОНТЕКСТУ, ОСТАЛЬНОЕ РАБОТАЕТ =================
	t.Nil(err)

}
