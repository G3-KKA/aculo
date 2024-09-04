package service

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	repository "aculo/batch-inserter/internal/repo"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"
	"sync"
	"sync/atomic"
)

// Static check
func _() {

	var _ Service = (*service)(nil)
	var _, _, _ = unifaces.Tx[ServiceAPI]((*service)(nil)).Tx()
}

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	unifaces.Tx[ServiceAPI]
	ServiceAPI
}

//go:generate mockery --filename=mock_service_api.go --name=ServiceAPI --dir=. --structname MockServiceAPI  --inpackage=true
type ServiceAPI interface {
	SendBatch(ctx context.Context, batch []domain.Event) error
	GracefulShutdown() error
}
type service struct {
	repo   unifaces.Tx[repository.RepositoryAPI]
	logger logger.Logger

	unavailable atomic.Bool
	mx          *sync.RWMutex
}

// # Common middleware for all api calls
//
// Safe to call multiple times, will return [unifaces.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic
func (s *service) Tx() (ServiceAPI, unifaces.TxClose, error) {

	if s.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}
	s.mx.RLock()
	var closed atomic.Bool
	f := func() error {
		if !closed.CompareAndSwap(false, true) {
			return unifaces.ErrTxAlreadyClosed
		}
		s.mx.RUnlock()
		return nil
	}
	return s, unifaces.TxClose(f), nil

}

// GracefulShutdown implements
func (s *service) GracefulShutdown() error {
	if s.unavailable.CompareAndSwap(false, true) {
		// GracefulShutdown can be called only if Tx was called, so we need a bit of magic here
		s.mx.RUnlock()
		// Lock cant be acquired while Tx with
		s.mx.Lock()
		defer s.mx.RLock()
		defer s.mx.Unlock()
		service, txclose, err := s.repo.Tx()
		defer txclose()
		if err != nil {
			return err
		}
		return service.GracefulShutdown()
	}
	return unierrors.ErrUnavailable
}

func (s *service) SendBatch(ctx context.Context, eventbatch []domain.Event) error {
	api, txclose, err := s.repo.Tx()
	if err != nil {
		return err
	}
	defer txclose()
	return api.SendBatch(ctx, eventbatch)

}
func New(ctx context.Context, config config.Config, l logger.Logger, repo unifaces.Tx[repository.RepositoryAPI]) (*service, error) {

	srvc := &service{
		repo:   repo,
		logger: l,
		mx:     &sync.RWMutex{},
	}

	return srvc, nil
}
