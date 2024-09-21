package app

import (
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/logger"
	repository "aculo/batch-inserter/internal/repo"
	"aculo/batch-inserter/internal/service"
	"context"

	"golang.org/x/sync/errgroup"
)

type App struct {
	controller txface.Tx[controller.ControllerAPI]
}

// ИНКАПСУЛИРОВАТЬ HTTP СЕРВЕР  В КОНТРОЛЛЕРА, ИБО КОНТРОЛЛЕР ЭТО ТАКЖЕ GRPC СЕРВИС И ОСТАЛЬНОЕ
func Run() (err error) {

	initConfig, err := config.ReadInConfig()
	if err != nil {
		return err
	}
	config.Print(initConfig)

	logger, _, err := logger.AssembleLogger(initConfig)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	app, err := newapp(ctx, initConfig, logger)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	errorGroup := &errgroup.Group{}
	wrapper := func() error {
		return app.Serve(ctx)
	}
	errorGroup.Go(wrapper)
	err = errorGroup.Wait()
	if err != nil {
		logger.Fatal(err)
	}
	return

}

const TODO_REPLACE_WITH_MASTERNODE = "NOT_EXISTING_TOPIC"

func newapp(ctx context.Context, conf config.Config, logger logger.Logger) (*App, error) {

	repo, err := repository.New(ctx, conf, logger)
	if err != nil {
		return nil, err
	}
	service, err := service.New(ctx, conf, logger, repo)
	if err != nil {
		return nil, err
	}
	controller, err := controller.New(ctx, TODO_REPLACE_WITH_MASTERNODE, conf, logger, service)

	if err != nil {
		return nil, err
	}
	return &App{
		controller: controller,
	}, nil
}
func (app *App) Serve(ctx context.Context) (err error) {

	// HTTP innterface нужен здесь
	// принимаем новый запрос на /metadata, запускае новый ctrl на топик
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		f := func() error {
			ctrl, txclose, err := app.controller.Tx()
			if err != nil {
				return err
			}
			defer txclose()
			return ctrl.HandleBatch(ctx)
		}
		err = f()
		if err != nil {
			return err
		}
	}

}

// OFFSET PROBLEM
func (app *App) Serve2(ctx context.Context) (err error) {
	errgroup := errgroup.Group{}
	errgroup.SetLimit(10)
	f := func() error {
		ctrl, txclose, err := app.controller.Tx()
		if err != nil {
			return err
		}
		defer txclose()
		return ctrl.HandleBatch(ctx)
	}
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}

		errgroup.Go(f)
	}
	return errgroup.Wait()

}
