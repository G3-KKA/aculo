package service

import (
	"aculo/connector-restapi/internal/brocker"
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/request"
	"context"
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	SendEvent(ctx context.Context, req request.SendEventRequest) (request.SendEventResponse, error)
}
type eventService struct {
	b brocker.Brocker
}
type BuildServiceRequest struct {
	Brocker brocker.Brocker
}

func New(ctx context.Context, config config.Config, brocker brocker.Brocker) (*eventService, error) {

	return &eventService{
		b: brocker,
	}, nil
}

func (s *eventService) SendEvent(ctx context.Context, req request.SendEventRequest) (request.SendEventResponse, error) {

	_, err := s.b.SendEvent(ctx, request.SendEventRequest{
		Topic: req.Topic,
		Event: req.Event,
	})
	if err != nil {
		return request.SendEventResponse{}, err
	}
	return request.SendEventResponse{}, nil
}
