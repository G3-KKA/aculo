package service

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	repository "aculo/batch-inserter/internal/repo"
	"context"
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	SendBatch(context.Context, []domain.Event) error
}

type service struct {
	repo repository.Repository
}

func (s *service) SendBatch(ctx context.Context, eventbatch []domain.Event) error {
	return s.repo.SendBatch(ctx, eventbatch)
}
func New(ctx context.Context, config config.Config, repo repository.Repository) (Service, error) {
	return &service{
		repo: repo,
	}, nil
}
