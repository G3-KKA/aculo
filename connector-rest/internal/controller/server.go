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

type server struct {
	serverHTTP http.Server
}

func New(ctx context.Context, cfg config.Config, service service.Service) (*server, error) {
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
	log.Info("Serving on: ", srv.serverHTTP.Addr)
	return srv.serverHTTP.ListenAndServe()
}
