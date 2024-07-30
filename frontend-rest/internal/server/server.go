package server

import (
	"aculo/frontend-restapi/internal/config"
	log "aculo/frontend-restapi/internal/logger"
	"context"
	"net/http"
)

//go:generate mockery --filename=mock_server.go --name=Server --dir=. --structname MockServer  --inpackage=true
type Server interface {
	ListenAndServe() error
}
type frontServer struct {
	stdServer *http.Server
}

func New(ctx context.Context, config config.Config, mux AppMux) Server {
	server := http.Server{
		Addr:    config.HTTPServer.ListeningAddress + config.HTTPServer.Port,
		Handler: mux,
	}
	return &frontServer{stdServer: &server}
}

func (s *frontServer) ListenAndServe() error {
	log.Info("starting server on: ", s.stdServer.Addr)
	return s.stdServer.ListenAndServe()
}
