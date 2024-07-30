package service

import (
	"aculo/connector-restapi/internal/brocker"
	"aculo/connector-restapi/internal/config"
	"context"
)

type SendEventRequest struct {
	Topic string
	Event []byte
}
type SendEventResponse struct{}

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname MockService  --inpackage=true
type Service interface {
	SendEvent(context.Context, SendEventRequest) (SendEventResponse, error)
}

type eventService struct {
	b brocker.Brocker
}
type BuildServiceRequest struct {
	Brocker brocker.Brocker
}

func New(ctx context.Context, config config.Config, req BuildServiceRequest) (Service, error) {

	return &eventService{
		b: req.Brocker,
	}, nil
}

func (s *eventService) SendEvent(ctx context.Context, req SendEventRequest) (SendEventResponse, error) {

	_, err := s.b.SendEvent(ctx, brocker.SendEventRequest{
		Topic: req.Topic,
		Event: req.Event,
	})
	if err != nil {
		return SendEventResponse{}, err
	}
	return SendEventResponse{}, nil
}
