package service

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/controller/broker"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/logger"
	repository "aculo/batch-inserter/internal/repo"
	"aculo/batch-inserter/internal/unified/unierrors"
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Static check
func _() {
	var _ Service = (*service)(nil)
	var _, _, _ = txface.Tx[ServiceAPI]((*service)(nil)).Tx()
}

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	txface.Tx[ServiceAPI]
}

//go:generate mockery --filename=mock_service_api.go --name=ServiceAPI --dir=. --structname MockServiceAPI  --inpackage=true
type ServiceAPI interface {
	SendBatch(ctx context.Context, batch []domain.Log) error
	Shutdown() error
	HandleNewClient(ctx context.Context) (topic string, err error)
}
type txprovider struct {
	s      *service
	repo   repository.Repository
	broker broker.Broker
}
type service struct {
	//
	repo txface.Tx[repository.RepositoryAPI] //txface.Tx[repository.RepositoryAPI]
	//
	broker txface.Tx[broker.BrokerAPI]
	logger logger.Logger

	unavailable atomic.Bool
	mx          *sync.RWMutex
}

// HandleNewLogTopic implements ServiceAPI.
func (s *service) HandleNewClient(ctx context.Context) (topic string, err error) {

	//
	broker, btxclose, err := s.broker.Tx()
	if err != nil {
		return "", err
	}
	defer btxclose()

	// Get log channel
	topic, err = broker.NewTopic(ctx)
	if err != nil {
		return "", err
	}
	logs, err := broker.HandleTopic(ctx, topic)
	if err != nil {
		broker.DeleteTopic(ctx, topic)
		return "", err
	}

	//
	repo, rtxclose, err := s.repo.Tx()
	if err != nil {
		err2 := broker.StopHandling(ctx, topic)
		errors.Join(err2, err)
		err2 = broker.DeleteTopic(ctx, topic)
		errors.Join(err2, err)
		return "", err
	}
	defer rtxclose()

	// Send logs to repo
	err = repo.HandleLogStream(ctx, logs)
	if err != nil {
		err2 := broker.StopHandling(ctx, topic)
		errors.Join(err2, err)
		err2 = broker.DeleteTopic(ctx, topic)
		errors.Join(err2, err)
		return "", err
	}
	return topic, err
}

// # Common middleware for all api calls
//
// Safe to call [unifaces.TxClose] multiple times,
// will return [unifaces.ErrTxAlreadyClosed] not breaking logic.
func (s *service) Tx() (ServiceAPI, txface.Commit, error) {

	if s.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	//
	// Actually start the transaction
	s.mx.RLock()
	// TODO: if someone sets unavaliable = true AND acquire Write Lock()
	// Faster than we take RLock -- we could  access memory that in process of shutting down
	if s.unavailable.Load() {
		s.mx.RUnlock()
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	//
	// Multiple txclose() call safe measure
	var closed atomic.Bool
	txclose := func() error {

		if !closed.CompareAndSwap(false, true) {
			return txface.ErrTxAlreadyClosed
		}

		s.mx.RUnlock()
		return nil
	}

	return s, txface.Commit(txclose), nil

}

// GracefulShutdown implements
func (s *txprovider) Shutdown() error {
	if s.unavailable.CompareAndSwap(false, true) {
		s.mx.Lock()
		defer s.mx.Unlock()
		return s.repo.S
	}
	return unierrors.ErrUnavailable
}

type Metadata struct {
	topic   string
	address string
}

/*
func (s *service) HandleNewTopic(ctx context.Context) (topic string, err error) {

	broker, txclose, err := s.broker.Tx()
	if err != nil {
		return "", err
	}
	defer txclose()

	topic, err = broker.NewTopic(ctx)
	if err != nil {
		return "", err
	}
	logs, err := broker.HandleTopic(ctx, topic)
	if err != nil {
		err2 := broker.DeleteTopic(ctx, topic)
		err = errors.Join(err, err2)
		return "", err
	}
	return

}
*/
func (s *service) SendBatch(ctx context.Context, eventbatch []domain.Log) error {
	repoapi, txclose, err := s.repo.Tx()
	if err != nil {
		return err
	}
	defer txclose()
	return repoapi.SendBatch(ctx, eventbatch)

}
func New(
	ctx context.Context,
	config config.Config,
	l logger.Logger,
	repo txface.Tx[repository.RepositoryAPI],
	broker txface.Tx[broker.BrokerAPI],
) (*service, error) {

	srvc := &service{
		repo:   repo,
		logger: l,
		mx:     &sync.RWMutex{},
	}

	return srvc, nil
}
