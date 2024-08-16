package repository

import (
	"aculo/frontend-restapi/domain"
	"aculo/frontend-restapi/internal/config"
	log "aculo/frontend-restapi/internal/logger"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type GetEventRequest struct {
	EID string
}
type GetEventResponse struct {
	Event domain.Event
}

//go:generate mockery --filename=mock_repository.go --name=Repository --dir=. --structname MockRepository  --inpackage=true
type Repository interface {
	GetEvent(ctx context.Context, req GetEventRequest) (GetEventResponse, error)
}

func New(ctx context.Context, config config.Config) (Repository, error) {
	conn, closeconn, err := ErrorproofGetConnect()
	if err != nil {
		return nil, err
	}
	go func() {
		//  TODO ===========  close() нормальным способом ===========
		<-ctx.Done()
		err := closeconn()
		if err != nil {
			log.Info("close error: ", err)
		}
	}()
	repo := &eRepo{conn: conn}
	return repo, nil

}

type eRepo struct {
	conn clickhouse.Conn
}

// GetEvent implements EventRepository.
func (e *eRepo) GetEvent(ctx context.Context, req GetEventRequest) (GetEventResponse, error) {
	chCtx := clickhouse.Context(context.TODO(),
		clickhouse.WithParameters(clickhouse.Parameters{
			"eid": req.EID,
		}))

	row := e.conn.QueryRow(chCtx, "SELECT * FROM event.main_table WHERE eid = {eid:String} LIMIT 1")

	event := domain.Event{}
	if err := row.ScanStruct(&event); err != nil {
		return GetEventResponse{}, err
	}
	return GetEventResponse{
		Event: event,
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

func ErrorproofGetConnect() (driver.Conn, ClickhouseCloseFunc, error) {
	conn, err := Click()
	return conn, ClickhouseCloseFunc(conn.Close), err
}
