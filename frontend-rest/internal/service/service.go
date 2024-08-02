package service

import (
	"aculo/frontend-restapi/internal/config"
	log "aculo/frontend-restapi/internal/logger"
	repository "aculo/frontend-restapi/internal/repo"
	eventservice "aculo/frontend-restapi/internal/service/event"
	"aculo/frontend-restapi/internal/service/transfomer"
	"context"
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	eventservice.EventService
}
type service struct {
	eventservice.EventService
}

type BuildServiceRequest struct {
	Repo repository.Repository
}

func New(ctx context.Context, config config.Config, repo repository.Repository) (Service, error) {

	eservice, err := eventservice.New(ctx, config, eventservice.BuildEserviceRequest{
		Repo:        repo,
		Transformer: transfomer.New(ctx, config),
	})
	if err != nil {
		log.Info("assemble event service failed: ", err)
		return nil, err
	}
	return &service{
		EventService: eservice,
	}, nil

}
