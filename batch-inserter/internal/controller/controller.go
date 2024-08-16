package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

var _ unifaces.Tx[ControllerAPI] = (*controller)(nil)
var _ Controller = (*controller)(nil)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
type Controller interface {
	Tx() (ControllerAPI, unifaces.TxClose, error)
}

//go:generate mockery --filename=mock_controllerapi.go --name=ControllerAPI --dir=. --structname MockControllerAPI  --inpackage=true
type ControllerAPI interface {
	HandleBatch(ctx context.Context) error
	GracefulShutdown() error
}

type controller struct {
	service       service.Service
	config        config.Config
	batchprovider *BatchProvider[domain.Event]
	logger        logger.Logger

	unavailable atomic.Bool
	//txcounter   atomic.Int64
	mx *sync.RWMutex
}

func New(ctx context.Context, config config.Config, l logger.Logger, service service.Service) (Controller, error) {

	return &controller{
		service:       service,
		config:        config,
		batchprovider: NewBatchProvider[domain.Event](ctx, config),
		logger:        l,

		unavailable: atomic.Bool{},
		// txcounter:   atomic.Int64{}, // TODO , можно обойтись без счетчика, его роль исполнит мьютекс, там этот счётчик уже есть
		mx: &sync.RWMutex{},
	}, nil
}

// # Common middleware for all api calls
//
// Safe to close multiple times, will return [unifaces.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic
//
// Safe to close if error happened, but not necessary,
// In case of error [unifaces.TxClose] does nothing but return same error
func (ctl *controller) Tx() (self ControllerAPI, txclose unifaces.TxClose, err error) {

	if ctl.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	ctl.mx.RLock()
	var closed atomic.Bool
	f := func() error {
		if !closed.CompareAndSwap(false, true) {
			return unifaces.ErrTxAlreadyClosed
		}
		ctl.mx.RUnlock()
		return nil
	}
	return ctl, unifaces.TxClose(f), nil

}

// GracefulShutdown implements RepositoryAPI.
func (ctl *controller) GracefulShutdown() error {
	if ctl.unavailable.CompareAndSwap(false, true) {
		// GracefulShutdown can be called only if Tx was called, so we need a bit of magic here
		ctl.mx.RUnlock()
		// Lock cant be acquired while Tx with
		ctl.mx.Lock()
		defer ctl.mx.RLock()
		defer ctl.mx.Unlock()

		service, txclose, err := ctl.service.Tx()
		defer txclose()

		if err != nil {
			return err
		}
		return service.GracefulShutdown()
	}
	return unierrors.ErrUnavailable
}

// TODO , use ctx to cancel in the middle
// Returns [unierrors.ErrOperationInterrupted] if context was canceled in the middle
func (ctl *controller) HandleBatch(ctx context.Context) error {
	// flowController.RegisterChildren(service) (children ...Registrable)
	cfg := ctl.config.Broker

	consumer, close, err := assembleConsumer(cfg)
	if err != nil {
		return err
	}
	defer close()
	ch := consumer.Messages()

	events, returnbatch := ctl.batchprovider.GetBatch()
	defer returnbatch()
	var msg *sarama.ConsumerMessage
	for i := range cfg.BatchSize {
		select {
		case <-ctx.Done():
			return unierrors.ErrOperationInterrupted
		case msg = <-ch: // common case
		}

		// TODO msg.Headers
		//
		// ctl.logger.Debug("TODO handwriting event, make it make sence, idiot, use message headers ")
		events[i].Data = msg.Value
		events[i].EID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
		events[i].ProviderID = "0"
		events[i].SchemaID = "0"
		events[i].Type = "test_type"
		// ctl.logger.Debug("Event received", string(events[i].Data))

	}
	service, txclose, err := ctl.service.Tx()
	if err != nil {
		return err
	}
	defer txclose()
	ctl.logger.Debug("sending batch")
	err = service.SendBatch(ctx, events)
	if err != nil {
		ctl.logger.Info("failed to send batch: %v", err)
	}
	return err

}

type mustDeferFunc func() error

func assembleConsumer(cfg config.Broker) (sarama.PartitionConsumer, mustDeferFunc, error) {

	consumer, err := sarama.NewConsumer(cfg.Addresses, nil)
	if err != nil {
		return nil, nil, err
	}
	// TODO hardcode here !!
	partitionConsumer, err := consumer.ConsumePartition(cfg.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		consumer.Close()
		return nil, nil, err
	}
	consumerClose := func() error {
		err := partitionConsumer.Close()
		if err != nil {
			return ErrFailedToCloseConsumer
		}
		err = consumer.Close()
		if err != nil {
			return ErrFailedToCloseConsumer
		}
		return nil
	}
	return partitionConsumer, consumerClose, nil
}
