package app

import (
	"aculo/frontend-restapi/internal/config"
	log "aculo/frontend-restapi/internal/logger"
	repository "aculo/frontend-restapi/internal/repo"
	"aculo/frontend-restapi/internal/server"
	"aculo/frontend-restapi/internal/server/groups/event"
	"aculo/frontend-restapi/internal/service"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type App struct {
	config config.Config
	server server.Server

	mux server.AppMux
}

// May be useful for
func (app *App) Mux() server.AppMux {
	return app.mux
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
	app, err := new(context.TODO(), initConfig)
	if err != nil {
		return err
	}

	// Running server
	errorGroup := &errgroup.Group{}
	errorGroup.Go(app.server.ListenAndServe)

	err = errorGroup.Wait()
	return
}

func new(ctx context.Context, config config.Config) (app *App, err error) {
	conn, close, err := ErrorproneGetConnect()
	if err != nil {
		return nil, err
	}
	defer close()

	// Assembling repository
	repo, err := repository.New(ctx, config, conn)
	if err != nil {
		log.Info("assemble event repository failed: ", err)

		return nil, err
	}

	// Assembling service
	service, err := service.New(ctx, config, service.BuildServiceRequest{
		Repo: repo,
	})
	if err != nil {
		log.Info("assemble service failed: ", err)
		return nil, err
	}

	// Preparing endpoints
	rootEndpoints := []server.Attachable{
		event.NewSpecialGroup(),
	}

	chains := []server.Chain{
		chain(event.NewEventGroup(ctx, config, service), event.NewSpecialGroup()),
	}

	root := gin.New()

	mux := server.NewMux(ctx, config, root, rootEndpoints, chains)
	server := server.New(ctx, config, mux)
	return &App{
		server: server,
		mux:    mux,
		config: config,
	}, nil
}

// TODO make normal

func Click() (driver.Conn, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"}, // TODO get from config, stop hardcode
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return conn, conn.Ping(context.Background())

}

// Errorprone Get conn
type ClickhouseCloseFunc func() error

func ErrorproneGetConnect() (driver.Conn, ClickhouseCloseFunc, error) {
	conn, err := Click()
	return conn, ClickhouseCloseFunc(conn.Close), err
}
