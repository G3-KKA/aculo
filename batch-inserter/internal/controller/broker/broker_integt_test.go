package broker

import (
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/testing/asyncsuite"
	"aculo/batch-inserter/internal/testing/testmark"
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

const (
	TEST_DEADLINE = time.Second * 10
)

type BrokerTestSuite_Intgr struct {
	suite.Suite
	logger          logger.Logger
	cfg             config.Config
	continueChannel chan os.Signal
}

func TestBroker(t *testing.T) {
	s := new(BrokerTestSuite_Intgr)

	testmark.MarkAs(testmark.INTEGRATION_TEST, s)

	suite.Run(t, s)
}

func (t *BrokerTestSuite_Intgr) SetupSuite() {

	//
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

	// Print Config
	config.Print(cfg)
	t.logger.Info("integr_test_config:", cfg)

	// Do not tear down untill we are done
	// Creates .continue file down the test
	// TODO: wrap it , too ugly
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

func (t *BrokerTestSuite_Intgr) BeforeTest(suiteName, testName string) {
	os.Stdout.WriteString("[INFO] " + testName + " starting \n")
	switch testName {
	case "Test_Integr_StreamLogs":
	default:
	}

}

func (t *BrokerTestSuite_Intgr) AfterTest(suiteName, testName string) {
	os.Stdout.WriteString("[INFO] " + testName + " cleanup \n")
	switch testName {
	case "Test_Integr_StreamLogs":
	}
}

func (t *BrokerTestSuite_Intgr) TearDownSuite() {
}

func (t *BrokerTestSuite_Intgr) Test_Integr_StreamLogs() {

	var (
		err     error
		topic   string
		reqdata []byte
	)
	reqdata = []byte(`{"id":1, "name":"test"}`)

	ctx, cancel := context.WithCancel(context.Background())

	btx, err := New(ctx, t.cfg, t.logger)
	t.NoError(err)

	//
	// We also want to test txclose() mechanism
	// So every interaction with broker will happen
	// In its own transaction
	newtopic := func() {

		//
		// Aquiring broker api  via Tx()
		broker, txclose, err := btx.Tx()
		t.NoError(err)
		defer txclose()

		//
		// Create topic and save its name
		topic, err = broker.NewTopic(ctx)
		t.NoError(err)
		t.logger.Info("topic are ", topic)
	}
	newtopic()

	producer, err := sarama.NewSyncProducer(t.cfg.Broker.Addresses, nil)
	t.NoError(err)
	asyncT := asyncsuite.AsyncSuite(&t.Suite)
	go func() {

		defer producer.Close()

		for {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(reqdata),
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			_, _, err := producer.SendMessage(msg)
			asyncT.NoError(err)
		}
	}()

	// Preparation and time limit
	deadline := TEST_DEADLINE
	closechan := make(chan struct{})
	brokerchan := make(chan struct{})
	timer := time.NewTimer(deadline)

	// Closes the controller second before deadline
	go func() {
		defer func() {
			closechan <- struct{}{}
		}()
		// Everything should be shutted down faster than one second
		time.Sleep(deadline - time.Second)
		cancel()             // Interrupting the controller
		err = btx.Shutdown() // Shutting down the controller
		asyncT.NoError(err)

	}()
	go func() {
		defer func() {
			brokerchan <- struct{}{}
		}()

		//
		//
		broker, txclose, err := btx.Tx()
		asyncT.NoError(err)
		logs, err := broker.HandleTopic(ctx, topic)
		asyncT.NoError(err)
		txclose()
		//
		//
		counter := 0
		for log := range logs {
			if (counter % 50) == 0 {
				asyncT.T().Log(log)
			}
		}
	}()

	// Wait and check
	select {
	case <-timer.C:
		asyncT.FailNow("Timeout")
		return
	case <-closechan:
	}
	select {
	case <-timer.C:
		asyncT.FailNow("Timeout")
		return
	case <-brokerchan:
	}
	/*
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
	*/
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
