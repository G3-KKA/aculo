package app

import (
	"context"
	"master-service/internal/config"
	"master-service/internal/controller"
	mock_controller "master-service/internal/controller/mocks"
	"master-service/internal/logger"
	"os"
	"testing"
)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname=MockController --outpkg=mock_app
type Controller interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
type App struct {
	ctl Controller
	l   logger.Logger
}

// New.Run() wrapper
func Run() {
	app, err := New()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	err = app.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
func New() (*App, error) {
	var (
		err            error
		cfg            config.Config
		internalLogger logger.Logger
	)
	app := &App{}

	cfg, err = config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	internalLogger, err = logger.New(cfg.L)
	if err != nil {
		return nil, err
	}
	app.l = internalLogger
	ctl, err := controller.New(cfg.C, internalLogger, mock_controller.NewMockService(&testing.T{}))
	if err != nil {
		internalLogger.Error(err.Error())
		return nil, err
	}
	app.ctl = ctl

	internalLogger.Debug("app construction succeded")
	return app, nil

}
func (app *App) Run() error {
	err := app.ctl.Serve(context.TODO())
	if err != nil {
		app.l.Debug(err.Error())
	}
	return err
}
