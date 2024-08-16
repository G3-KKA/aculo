package eventservice

import (
	"aculo/frontend-restapi/domain"
	"aculo/frontend-restapi/internal/config"
	log "aculo/frontend-restapi/internal/logger"
	repository "aculo/frontend-restapi/internal/repo"
	"aculo/frontend-restapi/internal/service/transfomer"
	"context"
)

type GetEventRequest struct {
	EID string
}
type GetEventResponse struct {
	Event domain.Event
}

//go:generate mockery --name=EventService --dir=. --outpkg=mock_event_service --filename=mock_event_service.go --output=./mocks/event_service --structname MockEventService
type EventService interface {
	GetEvent(ctx context.Context, req GetEventRequest) (GetEventResponse, error)
}

type eventService struct {
	repo        repository.Repository
	transformer transfomer.Transformer
}
type BuildEserviceRequest struct {
	Repo        repository.Repository
	Transformer transfomer.Transformer
}

func New(ctx context.Context, config config.Config, req BuildEserviceRequest) (EventService, error) {

	return &eventService{
		repo:        req.Repo,
		transformer: req.Transformer,
	}, nil
}

func (s *eventService) GetEvent(ctx context.Context, req GetEventRequest) (GetEventResponse, error) {

	repoRsp, err := s.repo.GetEvent(ctx, repository.GetEventRequest{
		EID: req.EID,
	})
	if err != nil {
		return GetEventResponse{}, err
	}
	trsp, err := s.transformer.Transform(ctx, transfomer.TransformRequest{
		SpecifiedSchema: struct{}{},
		Data:            repoRsp.Event.Data,
	})
	log.Info("trsp: NOT IMPLEMENTED YET", trsp)
	if err != nil {
		return GetEventResponse{}, err
	}
	return GetEventResponse{
		Event: repoRsp.Event,
	}, nil
}
