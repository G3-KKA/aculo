package app

import (
	"master-service/internal/config"
	"os"
)

type (
	App struct {
	}
)

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
		err error
		cfg config.Config
		internalLogger 
	)

	cfg, err = config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	internalLogger , err = 23_SPT
	if err != nil {
		return nil, err
	}

	panic("")

}
func (app *App) Run() error {

}
