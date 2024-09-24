package app

import (
	"master-service/internal/config"
	"master-service/internal/logger"
	"os"
)

type (
	//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
	Controller interface {
		Serve() error
		Close() error
	}
	App struct {
		c Controller
	}
)

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

	cfg, err = config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	internalLogger, err = logger.New(cfg.L)
	if err != nil {
		return nil, err
	}

	panic("")

}
func (app *App) Run() error {

}
