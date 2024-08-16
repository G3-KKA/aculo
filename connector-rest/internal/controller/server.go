package controller

import (
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/controller/groups"
	"aculo/connector-restapi/internal/controller/groups/event"
	swaggergroup "aculo/connector-restapi/internal/controller/groups/swagger"
	log "aculo/connector-restapi/internal/logger"
	"aculo/connector-restapi/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
type Controller interface {
	Serve() error
}
type server struct {
	HTTPServer http.Server
}

func New(ctx context.Context, cfg config.Config, service service.Service) (Controller, error) {
	rootEndpoints := []groups.Attachable{
		swaggergroup.NewSwaggerGroup(),
	}

	chains := []groups.Chain{
		groups.MakeChain(event.NewEventGroup(ctx, cfg, service)),
	}

	engine := gin.New()

	mux := newMux(ctx, cfg, engine, rootEndpoints, chains)
	srv := server{
		http.Server{
			Addr:    config.AssembleAddress(cfg),
			Handler: mux,
		},
	}
	return &srv, nil
}
func (srv *server) Serve() error {
	log.Info("Serving on: ", srv.HTTPServer.Addr)
	return srv.HTTPServer.ListenAndServe()
}
