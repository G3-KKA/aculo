package controller

import (
	"aculo/frontend-restapi/internal/config"
	"aculo/frontend-restapi/internal/controller/groups"
	"aculo/frontend-restapi/internal/controller/groups/event"
	swaggergroup "aculo/frontend-restapi/internal/controller/groups/swagger"
	log "aculo/frontend-restapi/internal/logger"
	"aculo/frontend-restapi/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// / ==================================== ПЕРЕДЕЛАТЬ ВСЁ ЗДЕСЬ ====================================
//
//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
type Controller interface {
	Serve() error
}

var _ Controller = (*frontServer)(nil)

type frontServer struct {
	stdServer *http.Server
}

func New(ctx context.Context, config config.Config, srvc service.Service) (*frontServer, error) {
	// Preparing endpoints
	rootEndpoints := []groups.Attachable{
		event.NewSpecialGroup(), swaggergroup.NewSwaggerGroup(),
	}

	chains := []groups.Chain{
		groups.MakeChain(event.NewEventGroup(ctx, config, srvc), event.NewSpecialGroup()),
	}

	root := gin.New()

	mux := NewMux(ctx, config, root, rootEndpoints, chains)
	server := http.Server{
		Addr:    config.HTTPServer.ListeningAddress + config.HTTPServer.Port,
		Handler: mux,
	}
	return &frontServer{stdServer: &server}, nil
}

func (s *frontServer) Serve() error {
	log.Info("starting server on: ", s.stdServer.Addr)
	return s.stdServer.ListenAndServe()
}
