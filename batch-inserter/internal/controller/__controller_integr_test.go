package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"aculo/batch-inserter/internal/testutils"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite_Intgr struct {
	suite.Suite
	logger          logger.Logger
	admin           sarama.ClusterAdmin
	cfg             config.Config
	continueChannel chan os.Signal
}

func TestController(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite_Intgr))
}

const INTER_TEST_TOPIC = "intertesttopic"

func (t *ControllerTestSuite_Intgr) SetupSuite() {

	testutils.ThisIsIntegrationTest(t) // Mark this test as integration test.

	// Config
	cfg, err := config.ReadInConfig()
	t.NoError(err)
	t.NotZero(cfg)
	t.cfg = cfg

	// Logger
	mock_logger := logger.NewMockLogger(t.T())
	mock_logger.On("Info", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		t.T().Log(args)
	})
	mock_logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		t.T().Log(args)
	})
	t.logger = mock_logger
	config.Print(cfg)
	t.logger.Info("integr_test_config:", cfg)

	// Admin
	Admin, err := sarama.NewClusterAdmin(cfg.Broker.Addresses, nil)
	t.NoError(err)
	t.admin = Admin

	// Do not tear down untill we are done
	continuetesting := make(chan os.Signal, 1)
	signal.Notify(continuetesting, syscall.SIGUSR1)
	t.continueChannel = continuetesting

	_ = viper.BindEnv("DEBUG")
	if viper.Get("DEBUG") == nil {
		// Continue testing, if we do not need information
		// That will be lost after tear down
		t.continueChannel <- syscall.SIGUSR1
	}

}

func (t *ControllerTestSuite_Intgr) BeforeTest(suiteName, testName string) {
	os.Stdout.WriteString("[INFO] " + testName + " starting \n")
	switch testName {
	case "Test_Integr_HandleBatch":
		err := t.admin.CreateTopic(INTER_TEST_TOPIC, &sarama.TopicDetail{NumPartitions: 1, ReplicationFactor: 1}, false)
		t.NoError(err)
	default:
	}

}

func (t *ControllerTestSuite_Intgr) AfterTest(suiteName, testName string) {
	os.Stdout.WriteString("[INFO] " + testName + " cleanup \n")
	switch testName {
	case "Test_Integr_HandleBatch":
		err := t.admin.DeleteTopic(INTER_TEST_TOPIC)
		t.NoError(err)
	}
}

func (t *ControllerTestSuite_Intgr) TearDownSuite() {
}

func (t *ControllerTestSuite_Intgr) Test_Integr_HandleBatch() {

	reqdata := []byte(`{"id":1, "name":"test"}`)

	// Mock API
	mock_serviceapi := service.NewMockServiceAPI(t.T())
	mock_serviceapi.On("SendBatch", mock.Anything, mock.
		MatchedBy(func(batch []domain.Event) bool {
			if t.NotEmpty(batch) {
				t.Equal(batch[0].Data, reqdata)
			}
			return t.Equal(t.cfg.Broker.BatchSize, len(batch))
		})).Return(nil)

	// NOW IT HAPPENS IN masternode
	//mock_serviceapi.On("GracefulShutdown").Return(nil)

	// Mock Service
	mock_service := service.NewMockService(t.T())
	mock_service.On("Tx").Return(mock_serviceapi,
		unifaces.TxClose(func() error { return nil }),
		nil, // error nil
	)

	// Controller
	ctx, cancel := context.WithCancel(context.Background())
	ctrl, err := New(ctx,
		INTER_TEST_TOPIC,
		config.Config{
			Broker: t.cfg.Broker,
		},
		t.logger,
		mock_service,
	)
	t.NoError(err)

	// Real producer

	producer, err := sarama.NewSyncProducer(t.cfg.Broker.Addresses, nil)
	t.NoError(err)
	go func() {

		defer producer.Close()
		msg := &sarama.ProducerMessage{
			Topic: INTER_TEST_TOPIC,
			Value: sarama.ByteEncoder(reqdata),
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			_, _, err := producer.SendMessage(msg)
			t.NoError(err)
		}
	}()
	// sarama.DeleteTopicsRequest

	// Preparation and time limit
	deadline := 10 * time.Second
	closechan := make(chan struct{})
	controllerchan := make(chan struct{})
	timer := time.NewTimer(deadline)

	// Closes the controller second before deadline
	go func() {
		defer func() {
			closechan <- struct{}{}
		}()
		// Everything should be shutted down faster than one second
		time.Sleep(deadline - time.Second)
		api, txclose, err := ctrl.Tx()
		t.NoError(err)
		defer txclose()
		cancel()                     // Interrupting the controller
		err = api.GracefulShutdown() // Shutting down the controller
		t.NoError(err)

	}()
	go func() {
		defer func() {
			controllerchan <- struct{}{}
		}()
		for {
			api, txclose, err := ctrl.Tx()
			if err != nil {
				t.ErrorIs(err, unierrors.ErrUnavailable)
				return
			}
			t.NoError(err)
			err = func() error {
				defer txclose()
				return api.HandleBatch(ctx)
			}()
			if err != nil {
				t.ErrorIs(err, unierrors.ErrOperationInterrupted)
				return
			}
			t.NoError(err)
		}
	}()

	// Wait and check
	select {
	case <-timer.C:
		t.FailNow("Timeout")
		return
	case <-closechan:
	}
	select {
	case <-timer.C:
		t.FailNow("Timeout")
		return
	case <-controllerchan:
	}
	t.True(ctrl.unavailable.Load())
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	err = ctrl.HandleBatch(ctx)
	t.ErrorIs(err, unierrors.ErrOperationInterrupted)
	_, txclose, err := ctrl.Tx()
	t.ErrorIs(err, unierrors.ErrUnavailable)
	txclose()
	err = ctrl.GracefulShutdown()
	// Tx
	t.ErrorIs(err, unierrors.ErrUnavailable)
	// TODO: Hardcoded .continue path
	continueF, err := os.Create(viper.GetString("WORKSPACE") + "/tmp/.continue")
	t.NoError(err)

	// if DEBUG defined we wait for signal
	// otherwise it will(should) be skipped
	<-t.continueChannel

	continueF.Close()
	err = os.Remove(viper.GetString("WORKSPACE") + "/tmp/.continue")
	t.NoError(err)

}
