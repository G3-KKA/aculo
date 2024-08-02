package app

import (
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/controller"
	log "aculo/batch-inserter/internal/logger"
	repository "aculo/batch-inserter/internal/repo"
	"aculo/batch-inserter/internal/service"
	"context"

	"golang.org/x/sync/errgroup"
)

type App struct {
	controller controller.Controller
}

// ИНКАПСУЛИРОВАТЬ HTTP СЕРВЕР  В КОНТРОЛЛЕРА, ИБО КОНТРОЛЛЕР ЭТО ТАКЖЕ GRPC СЕРВИС И ОСТАЛЬНОЕ
func Run() (err error) {

	err = config.InitConfig()
	if err != nil {
		return err
	}
	initConfig := config.Get()

	err = log.InitGlobalLogger(initConfig)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	app, err := newapp(ctx, initConfig)
	if err != nil {
		return err
	}
	errorGroup := &errgroup.Group{}
	errorGroup.Go(app.controller.Serve)

	return errorGroup.Wait()

}
func newapp(ctx context.Context, conf config.Config) (*App, error) {
	repo, err := repository.New(ctx, conf)
	if err != nil {
		return nil, err
	}
	service, err := service.New(ctx, conf, repo)
	if err != nil {
		return nil, err
	}
	controller, err := controller.New(ctx, conf, service)

	if err != nil {
		return nil, err
	}
	return &App{
		controller: controller,
	}, nil
}
