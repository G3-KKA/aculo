package service

import (
	"aculo/batch-inserter/internal/brocker"
	"aculo/batch-inserter/internal/config"
	repository "aculo/batch-inserter/internal/repo"
	"context"
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	ConsumeAndLoad(context.Context) error
}
type BuildServiceRequest struct {
	Repo    repository.Repository
	Brocker brocker.Brocker
}
type TODORENAMEservice struct {
}

func (s *TODORENAMEservice) ConsumeAndLoad(context.Context) error {
	return nil
}
func NewService(ctx context.Context, config config.Config, req BuildServiceRequest) (Service, error) {
	return &TODORENAMEservice{}, nil
}
