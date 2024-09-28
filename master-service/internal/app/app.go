package app

import (
	"context"
	"os"
	"testing"

	"master-service/internal/config"
	"master-service/internal/controller"
	"master-service/internal/logger"

	mock_controller "master-service/internal/controller/mocks"
)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname=MockController --outpkg=mock_app
type Controller interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Application itself.
type App struct {
	ctl Controller
	l   logger.Logger
}

// Run is an app.New.Run() wrapper.
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

// # [App] constructor.
//
// Stages:
//
//   - Read config,
//   - Initialize internal logger,
//   - Get [Controller].
func New() (*App, error) {
	var (
		err            error
		cfg            config.Config
		internalLogger logger.Logger
	)

	cfg, err = config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	internalLogger, err = logger.New(cfg.L)
	if err != nil {
		return nil, err
	}

	cfg.Print(&internalLogger)
	ctl, err := controller.New(cfg.C, internalLogger, mock_controller.NewMockService(&testing.T{}))
	if err != nil {
		internalLogger.Error(err.Error())

		return nil, err
	}
	app := &App{
		ctl: ctl,
		l:   internalLogger,
	}
	internalLogger.Debug("app construction succeeded")

	return app, nil

}

// Run calls Serve on underlying [Controller].
//
// Warning! Blocks execution indefinitely!
func (app *App) Run() error {
	err := app.ctl.Serve(context.TODO())
	if err != nil {
		app.l.Debug(err.Error())
	}

	return err
}
