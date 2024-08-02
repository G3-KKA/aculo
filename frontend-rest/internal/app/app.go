package app

import (
	"aculo/frontend-restapi/internal/config"
	"aculo/frontend-restapi/internal/controller"
	log "aculo/frontend-restapi/internal/logger"
	repository "aculo/frontend-restapi/internal/repo"

	"aculo/frontend-restapi/internal/service"
	"context"

	"golang.org/x/sync/errgroup"
)

type App struct {
	controller controller.Controller
}

// Start point of app
func Run() (err error) {

	// Config initialization
	err = config.InitConfig()
	if err != nil {
		return err
	}
	initConfig := config.Get()

	// Logger initialization
	err = log.InitGlobalLogger(initConfig)
	if err != nil {
		return err
	}
	// TODO: use pretty for this log.Info("config: ", initConfig)

	// App initialization
	ctx := context.TODO()
	app, err := newapp(ctx, initConfig)
	if err != nil {
		return err
	}

	// Running server
	errorGroup := &errgroup.Group{}
	errorGroup.Go(app.controller.Serve)

	err = errorGroup.Wait()
	return
}

func newapp(ctx context.Context, config config.Config) (app *App, err error) {

	// Repository
	repo, err := repository.New(ctx, config)
	if err != nil {
		log.Info("assemble event repository failed: ", err)
		return nil, err
	}

	// Service
	service, err := service.New(ctx, config, repo)
	if err != nil {
		log.Info("assemble service failed: ", err)
		return nil, err
	}
	// Controller
	controller, err := controller.New(ctx, config, service)
	if err != nil {
		log.Info("assemble controller failed: ", err)
		return nil, err
	}

	return &App{
		controller: controller,
	}, nil
}
