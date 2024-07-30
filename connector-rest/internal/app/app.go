package app

import (
	"aculo/connector-restapi/internal/brocker"
	"aculo/connector-restapi/internal/config"
	log "aculo/connector-restapi/internal/logger"
	"aculo/connector-restapi/internal/server"
	"aculo/connector-restapi/internal/server/groups"
	"aculo/connector-restapi/internal/service"
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type App struct {
	mux    server.AppMux
	config config.Config
}

// Used for external Servers
func (app *App) Mux() server.AppMux {
	return app.mux
}
func Run() (err error) {

	// Config initialization
	err = config.InitConfig()
	if err != nil {
		return err
	}
	initConfig := config.Get()

	// Logger initialization
	err = log.InitGlobalLogger(config.Get())
	if err != nil {
		return err
	}
	// TODO: use pretty for this log.Info("config: ", initConfig)

	// Running App
	app, err := new(context.TODO(), initConfig)
	if err != nil {
		return err
	}
	errorGroup := &errgroup.Group{}
	errorGroup.Go(app.runServerless)

	err = errorGroup.Wait()
	return
}

// Default LinstenAndServe
func (app *App) runServerless() (err error) {

	address := config.AssembleAddress(app.config)
	log.Info("Starting Server on ", address)

	err = app.mux.Run(address)
	return
}

func new(ctx context.Context, config config.Config) (app *App, err error) {

	brocker, err := brocker.New(ctx, config, brocker.BuildBrockerRequest{})
	if err != nil {
		log.Info("assemble event brocker failed: ", err)
		return nil, err
	}
	service, err := service.New(ctx, config, service.BuildServiceRequest{
		Brocker: brocker,
	})
	if err != nil {
		log.Info("assemble event service failed: ", err)
		return nil, err
	}

	rootEndpoints := []server.Attachable{}

	chains := []server.Chain{
		chain(groups.NewEventGroup(ctx, config, service)),
	}

	engine := gin.New()

	mux := server.NewMux(ctx, config, engine, rootEndpoints, chains)
	return &App{
		mux:    mux,
		config: config,
	}, nil
}
