package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"
	"errors"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

var _ unifaces.Tx[ControllerAPI] = (*topicHandler)(nil)
var _ Controller = (*topicHandler)(nil)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
type Controller interface {
	unifaces.Tx[ControllerAPI]
}

//go:generate mockery --filename=mock_controllerapi.go --name=ControllerAPI --dir=. --structname MockControllerAPI  --inpackage=true
type ControllerAPI interface {
	HandleBatch(ctx context.Context) error
	HandleTopic(ctx context.Context) error
	GracefulShutdown() error
}

type topicHandler struct {
	service       unifaces.Tx[service.ServiceAPI]
	config        config.Config
	batchprovider *BatchProvider[domain.Event]
	logger        logger.Logger
	consumer      sarama.Consumer
	pconsumer     sarama.PartitionConsumer

	unavailable atomic.Bool
	//txcounter   atomic.Int64
	mx *sync.RWMutex
}

var _ unifaces.Tx[MasterAPI] = (*masternode)(nil)
var _ Master = (*masternode)(nil)

//go:generate mockery --filename=mock_master.go --name=Master --dir=. --structname MockMaster  --inpackage=true
type Master interface {
	unifaces.Tx[MasterAPI]
}

//go:generate mockery --filename=mock_masterAPI.go --name=MasterAPI --dir=. --structname MockMasterAPI  --inpackage=true
type MasterAPI interface {
	HandleTopic(ctx context.Context, topic string) error
}
type masternode struct {
	l logger.Logger

	cfg config.Config

	admin   sarama.ClusterAdmin
	nodeCtx context.Context
	nodes   []node
	srvc    unifaces.Tx[service.ServiceAPI]

	unavailable atomic.Bool

	mx sync.RWMutex
}
type node struct {
	thandler *topicHandler
	cancel   context.CancelFunc
	done     chan struct{}
}

// # Common middleware for all api calls
//
// Safe to close multiple times, will return [unifaces.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic
//
// Safe to close if error happened, but not necessary,
// In case of error [unifaces.TxClose] does nothing but return same error
func (master *masternode) Tx() (MasterAPI, unifaces.TxClose, error) {

	if master.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	master.mx.RLock()
	var closed atomic.Bool
	f := func() error {
		if !closed.CompareAndSwap(false, true) {
			return unifaces.ErrTxAlreadyClosed
		}
		master.mx.RUnlock()
		return nil
	}
	return master, unifaces.TxClose(f), nil

}

// GracefulShutdown
func (master *masternode) GracefulShutdown() error {
	if master.unavailable.CompareAndSwap(false, true) {
		// GracefulShutdown can be called only if Tx was called, so we need a bit of magic here
		master.mx.RUnlock() // our own transaction's rlock
		// Lock cant be acquired while Tx with
		master.mx.Lock()
		defer master.mx.RLock()
		defer master.mx.Unlock()

		// ВОТ ТУТ МЫ ДОЛЖНЫ УБИТЬ ВСЕХ tHandler
		wg := sync.WaitGroup{}
		wg.Add(len(master.nodes))
		for _, node := range master.nodes {
			go func() {
				defer wg.Done()
				node.cancel()
				<-node.done
			}()

		}
		wg.Wait()

		service, txclose, err := master.srvc.Tx()
		defer txclose()

		if err != nil {
			return err
		}
		return service.GracefulShutdown()
	}
	return unierrors.ErrUnavailable
}
func NewMaster(ctx context.Context, config config.Config, l logger.Logger, service unifaces.Tx[service.ServiceAPI]) (*masternode, error) {
	admin, err := sarama.NewClusterAdmin(config.Broker.Addresses, nil)
	if err != nil {
		return nil, err
	}

	// ADMIN.CLOSE() нигде не вызывается
	return &masternode{
		l:           l,
		cfg:         config,
		admin:       admin,
		nodeCtx:     ctx,
		nodes:       []node{},
		srvc:        service,
		unavailable: atomic.Bool{},
		mx:          sync.RWMutex{},
	}, nil
}

func (master *masternode) HandleTopic(ctx context.Context, topic string) error {

	err := master.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
		//ReplicaAssignment: map[int32][]int32{},
		//ConfigEntries:     map[string]*string{},
	}, false)
	if err != nil {
		return err
	}
	ctx2, cancel := context.WithCancel(context.Background())
	tHandler, err := New(context.WithoutCancel(ctx), topic, master.cfg, master.l, master.srvc)
	if err != nil {
		return err
	}

	node := node{
		thandler: tHandler,
		cancel:   cancel,
		done:     make(chan struct{}, 1),
	}
	go func() {
		replyDoneF := func() {
			node.done <- struct{}{}
			close(node.done)
		}
		defer replyDoneF()
		handleF := func() error {
			tAPI, txclose, err := tHandler.Tx()
			if err != nil {
				return err
			}
			defer txclose()
			return tAPI.HandleTopic(ctx2)
		}
		err := handleF()
		if err != nil {
			master.l.Error(err)
		}
		graceShutF := func() error {
			tAPI, txclose, err := tHandler.Tx()
			if err != nil {
				return err
			}
			defer txclose()
			return tAPI.GracefulShutdown()
		}
		err = graceShutF()
		if err != nil {
			master.l.Error(err)
		}

	}()
	// ЗДЕСЬ ОСТАНОВИЛСЯ !!!!
	master.nodes = append(master.nodes, node)
	return nil

}
func New(ctx context.Context, topic string, config config.Config, l logger.Logger, service unifaces.Tx[service.ServiceAPI]) (*topicHandler, error) {
	consumer, err := sarama.NewConsumer(config.Broker.Addresses, nil)
	if err != nil {
		for range noexport_RETRY_COUNT {
			time.Sleep(noexport_RETRY_DELAY)
			consumer, err = sarama.NewConsumer(config.Broker.Addresses, nil)
			if err == nil {
				break
			}
		}
		err := errors.Join(ErrRetryNotWorksConnectionRefused, err)
		return nil, err
	}
	// TODO hardcode here !!
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		consumer.Close()
		return nil, err
	}
	return &topicHandler{
		config: config,
		logger: l,

		service:       service,
		batchprovider: NewBatchProvider[domain.Event](ctx, config),
		consumer:      consumer,
		pconsumer:     partitionConsumer,

		unavailable: atomic.Bool{},
		mx:          &sync.RWMutex{},
	}, nil
}

// # Common middleware for all api calls
//
// Safe to close multiple times, will return [unifaces.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic
//
// Safe to close if error happened, but not necessary,
// In case of error [unifaces.TxClose] does nothing but return same error
func (ctl *topicHandler) Tx() (ControllerAPI, unifaces.TxClose, error) {

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

// # GracefulShutdown
//
// Closes consumers
func (ctl *topicHandler) GracefulShutdown() (err error) {
	if ctl.unavailable.CompareAndSwap(false, true) {
		// GracefulShutdown can be called only if Tx was called, so we need a bit of magic here
		ctl.mx.RUnlock() // our own transaction's rlock
		// Lock cant be acquired while Tx with
		ctl.mx.Lock()
		defer ctl.mx.RLock()
		defer ctl.mx.Unlock()

		if ctl.consumer != nil {
			err = ctl.pconsumer.Close()
			if err != nil {
				ctl.logger.Error(err)
			}
			err = errors.Join(err, ctl.consumer.Close())
			if err != nil {
				ctl.logger.Error(err)
			}
		}

		/* 		service, txclose, err := ctl.service.Tx()
		   		defer txclose() */

		/* 		if err != nil {
		   			return err
		   		}
		   		return service.GracefulShutdown() */
		return err
	}
	return unierrors.ErrUnavailable
}

// TODO , use ctx to cancel in the middle
// Returns [unierrors.ErrOperationInterrupted] if context was canceled in the middle
func (ctl *topicHandler) HandleBatch(ctx context.Context) error {
	// flowController.RegisterChildren(service) (children ...Registrable)
	cfg := ctl.config.Broker

	ch := ctl.pconsumer.Messages()

	events, returnbatch := ctl.batchprovider.GetBatch()
	defer returnbatch()
	//events := make([]domain.Event, ctl.config.BatchSize)
	var msg *sarama.ConsumerMessage
	os.Stderr.WriteString(strconv.Itoa(len(events)))
	for i := range cfg.BatchSize {
		select {
		case <-ctx.Done():
			return unierrors.ErrOperationInterrupted
		case msg = <-ch: // common case
		}
		// TODO msg.Headers
		//
		//ctl.logger.Debug("TODO handwriting event, make it make sence, idiot, use message headers ")
		events[i].Data = msg.Value
		events[i].EID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
		events[i].ProviderID = "0"
		events[i].SchemaID = "0"
		events[i].Type = "test_type"
		//ctl.logger.Debug("Event received", string(events[i].Data))

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
func (ctl *topicHandler) HandleTopic(ctx context.Context) error {
	// flowController.RegisterChildren(service) (children ...Registrable)
	cfg := ctl.config.Broker

	ch := ctl.pconsumer.Messages()
	var err error
LOOP:
	for {
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			break LOOP
		default:
			f := func() error {
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
					//ctl.logger.Debug("TODO handwriting event, make it make sence, idiot, use message headers ")
					events[i].Data = msg.Value
					events[i].EID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
					events[i].ProviderID = "0"
					events[i].SchemaID = "0"
					events[i].Type = "test_type"
					//ctl.logger.Debug("Event received", string(events[i].Data))

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
			err = f()
		}
	}
	return err

}

type mustDeferFunc func() error

const noexport_RETRY_COUNT = 5
const noexport_RETRY_DELAY = 500 * time.Millisecond

func assembleConsumerNonPartion(addresses []string) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer(addresses, nil)
	if err != nil {
		for range noexport_RETRY_COUNT {
			time.Sleep(noexport_RETRY_DELAY)
			consumer, err = sarama.NewConsumer(addresses, nil)
			if err == nil {
				break
			}
		}
		err := errors.Join(ErrRetryNotWorksConnectionRefused, err)
		return nil, err
	}
	// TODO hardcode here !!
	return consumer, nil

}

// DEPRECATED !!!!
func assembleConsumer(topic string, cfg config.Broker) (sarama.PartitionConsumer, mustDeferFunc, error) {

	consumer, err := sarama.NewConsumer(cfg.Addresses, nil)
	if err != nil {
		for range noexport_RETRY_COUNT {
			time.Sleep(noexport_RETRY_DELAY)
			consumer, err = sarama.NewConsumer(cfg.Addresses, nil)
			if err == nil {
				break
			}
		}
		err := errors.Join(ErrRetryNotWorksConnectionRefused, err)
		return nil, nil, err
	}
	// TODO hardcode here !!
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
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
