package controller

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/generics/batchprovider"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/wrapper/asyncerror"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

// Я нихуя не понимаю что здесь происходит, господи что за пиздец, но тут почти всё нужный код, нужно лишь зарефакторить

var _ txface.Tx[ControllerAPI] = (*topicHandler)(nil)
var _ Controller = (*topicHandler)(nil)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
type Controller interface {
	txface.Tx[ControllerAPI]
}

//go:generate mockery --filename=mock_controllerapi.go --name=ControllerAPI --dir=. --structname MockControllerAPI  --inpackage=true
type ControllerAPI interface {
	HandleBatch(ctx context.Context) error
	HandleTopic(ctx context.Context) error
	GracefulShutdown() error
}

type topicHandler struct {
	service       txface.Tx[service.ServiceAPI]
	config        config.Config
	batchprovider *batchprovider.BatchProvider[domain.Log]
	logger        logger.Logger
	consumer      sarama.Consumer
	pconsumer     sarama.PartitionConsumer

	unavailable atomic.Bool
	//txcounter   atomic.Int64
	mx *sync.RWMutex
}

var _ txface.Tx[MasterAPI] = (*masternode)(nil)
var _ Master = (*masternode)(nil)

//go:generate mockery --filename=mock_master.go --name=Master --dir=. --structname MockMaster  --inpackage=true
type Master interface {
	txface.Tx[MasterAPI]
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
	srvc    txface.Tx[service.ServiceAPI]

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
// Safe to close multiple times, will return [txface.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic
//
// Safe to close if error happened, but not necessary,
// In case of error [txface.Commit] does nothing but return same error
func (master *masternode) Tx() (MasterAPI, txface.Commit, error) {

	if master.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	master.mx.RLock()
	// TODO: if someone sets unavaliable = true AND acquire Write Lock()
	// Faster than we take RLock -- we could  access memory that in process of shutting down
	if master.unavailable.Load() {
		master.mx.RUnlock()
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}
	var closed atomic.Bool
	f := func() error {
		if !closed.CompareAndSwap(false, true) {
			return txface.ErrTxAlreadyClosed
		}
		master.mx.RUnlock()
		return nil
	}
	return master, txface.Commit(f), nil

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

		service, Commit, err := master.srvc.Tx()
		defer Commit()

		if err != nil {
			return err
		}
		return service.GracefulShutdown()
	}
	return unierrors.ErrUnavailable
}
func NewMaster(ctx context.Context, config config.Config, l logger.Logger, service txface.Tx[service.ServiceAPI]) (*masternode, error) {
	admin, err := sarama.NewClusterAdmin(config.Broker.Addresses, nil)
	if err != nil {
		return nil, err
	}

	// ADMIN.CLOSE() нигде не вызывается
	return &masternode{
		l:       l,
		cfg:     config,
		admin:   admin,
		nodeCtx: ctx,
		nodes:   []node{},
		srvc:    service,

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
			tAPI, Commit, err := tHandler.Tx()
			if err != nil {
				return err
			}
			defer Commit()
			return tAPI.HandleTopic(ctx2)
		}
		err := handleF()
		if err != nil {
			master.l.Error(err)
		}
		graceShutF := func() error {
			tAPI, Commit, err := tHandler.Tx()
			if err != nil {
				return err
			}
			defer Commit()
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
func New(ctx context.Context, topic string, config config.Config, l logger.Logger, service txface.Tx[service.ServiceAPI]) (*topicHandler, error) {

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
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		consumer.Close()
		return nil, err
	}
	return &topicHandler{
		config: config,
		logger: l,

		service:       service,
		batchprovider: batchprovider.NewBatchProvider[domain.Log](ctx, config),
		consumer:      consumer,
		pconsumer:     partitionConsumer,

		unavailable: atomic.Bool{},
		mx:          &sync.RWMutex{},
	}, nil
}

// # Common middleware for all api calls
//
// Safe to close multiple times, will return [txface.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic
//
// Safe to close if error happened, but not necessary,
// In case of error [txface.Commit] does nothing but return same error
func (ctl *topicHandler) Tx() (ControllerAPI, txface.Commit, error) {

	if ctl.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	ctl.mx.RLock()
	var closed atomic.Bool
	f := func() error {
		if !closed.CompareAndSwap(false, true) {
			return txface.ErrTxAlreadyClosed
		}
		ctl.mx.RUnlock()
		return nil
	}
	return ctl, txface.Commit(f), nil

}

// # GracefulShutdown
//
// Closes consumers
func (ctl *topicHandler) GracefulShutdown() (err error) {
	if ctl.unavailable.CompareAndSwap(false, true) {

		//
		// GracefulShutdown can be called only if Tx was called, so we need a bit of magic here
		ctl.mx.RUnlock() // un'do our own transaction's rlock
		ctl.mx.Lock()    // and take fulllock

		//
		// re'do out own transaction's rlock, client
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

		/* 		service, Commit, err := ctl.service.Tx()
		   		defer Commit() */

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

	msgs := ctl.pconsumer.Messages()
	//  ?????
	events, returnbatch := ctl.batchprovider.GetBatch()
	defer returnbatch()
	//events := make([]domain.Event, ctl.config.BatchSize)
	var msg *sarama.ConsumerMessage
	//os.Stderr.WriteString(strconv.Itoa(len(events)))
	event := domain.Log{}
	for i := range cfg.BatchSize {
		select { // there are
		case <-ctx.Done():
			return unierrors.ErrOperationInterrupted
		case msg = <-msgs: // common case
		}
		// TODO msg.Headers
		//
		//ctl.logger.Debug("TODO handwriting event, make it make sence, idiot, use message headers ")
		event.Data = msg.Value
		// ОЧЕНЬ ДОРОГО !
		event.LogID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
		event.ProviderID = "0"
		event.SchemaID = "0"
		event.Type = "test_type"
		events[i] = event
		//ctl.logger.Debug("Event received", string(events[i].Data))

	}
	service, Commit, err := ctl.service.Tx()
	if err != nil {
		return err
	}
	defer Commit()
	//ctl.logger.Debug("sending batch")
	// go func here ?
	go func() {
		err = service.SendBatch(ctx, events)
		if err != nil {
			ctl.logger.Info("failed to send batch: %v", err)
		}
	}()
	return err

}
func (ctl *topicHandler) HandleTopic(ctx context.Context) error {

	cfg := ctl.config.Broker

	msgs := ctl.pconsumer.Messages()

	asyncerr := asyncerror.AsyncError{}
	for {
		err := ctl.HandleBatch(ctx)
		if err != nil {
			break
		}
	}
LOOP:
	for {
		if err2 := asyncerr.Err(); err2 != nil {
			return err2
		}
		select {
		case <-ctx.Done():
			break LOOP
		default:
			f := func() error {
				events, returnbatch := ctl.batchprovider.GetBatch()
				panichappened := atomic.Bool{}
				safemeasureF := func() {
					// if panic not happened
					if !panichappened.Load() {
						return
					}
					returnbatch()
				}
				defer safemeasureF()
				var msg *sarama.ConsumerMessage
				for i := range cfg.BatchSize {
					select {
					case <-ctx.Done():
						return unierrors.ErrOperationInterrupted
					case msg = <-msgs: // common case
					}
					// TODO msg.Headers
					//
					//ctl.logger.Debug("TODO handwriting event, make it make sence, idiot, use message headers ")
					events[i].Data = msg.Value
					events[i].LogID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
					events[i].ProviderID = "0"
					events[i].SchemaID = "0"
					events[i].Type = "test_type"
					//ctl.logger.Debug("Event received", string(events[i].Data))

				}
				service, Commit, err := ctl.service.Tx()
				if err != nil {
					return err
				}
				defer Commit()
				ctl.logger.Debug("sending batch")
				go func() {
					defer returnbatch() // unsafe return batch !!!, panic may happen before this G even created
					err := service.SendBatch(ctx, events)
					asyncerr.Join(err)

				}()
				return err
			}
			err := f()
			asyncerr.Join(err)
		}
	}
	return asyncerr.Err()

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
