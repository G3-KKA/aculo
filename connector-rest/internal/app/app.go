package app

import (
	"aculo/connector-restapi/internal/brocker"
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/controller"
	log "aculo/connector-restapi/internal/logger"
	"aculo/connector-restapi/internal/service"
	"context"

	"golang.org/x/sync/errgroup"
)

type App struct {
	controller controller.Controller
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
	ctx := context.TODO()

	// Running App
	app, err := newapp(ctx, initConfig)
	if err != nil {
		return err
	}

	errorGroup := &errgroup.Group{}
	errorGroup.Go(app.controller.Serve)

	return errorGroup.Wait()

}

func newapp(ctx context.Context, config config.Config) (app *App, err error) {

	brocker, err := brocker.New(ctx, config)
	if err != nil {
		log.Info("assemble brocker failed: ", err)
		return nil, err
	}
	service, err := service.New(ctx, config, brocker)
	if err != nil {
		log.Info("assemble service failed: ", err)
		return nil, err
	}
	controller, err := controller.New(ctx, config, service)
	if err != nil {
		log.Info("assemble controller failed: ", err)
		return nil, err
	}

	return &App{
		controller: controller,
	}, nil
}
