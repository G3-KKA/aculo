package server

import (
	"aculo/connector-restapi/internal/config"
	"context"
	"net/http"
)

type Server interface {
	ListenAndServe() error
}
type server struct {
	http.Server
}

func New(ctx context.Context, cfg config.Config, mux AppMux) Server {
	srv := server{
		http.Server{
			Addr:    config.AssembleAddress(cfg),
			Handler: mux,
		},
	}
	return &srv
}
