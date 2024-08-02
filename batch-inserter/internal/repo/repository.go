package repository

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	log "aculo/batch-inserter/internal/logger"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

//go:generate mockery --filename=mock_repository.go --name=Repository --dir=. --structname MockRepository  --inpackage=true
type Repository interface {
	SendBatch(context.Context, []domain.Event) error
}
type clickhouseRepo struct {
	ch clickhouse.Conn
}

func New(ctx context.Context, config config.Config) (Repository, error) {
	conn, err := Click()
	if err != nil {
		return nil, err
	}
	return &clickhouseRepo{
		ch: conn,
	}, nil

}
func (c *clickhouseRepo) SendBatch(ctx context.Context, eventbatch []domain.Event) error {
	batch, err := c.ch.PrepareBatch(ctx, "INSERT INTO event.main_table (eid, provider_id, schema_id, type, data)")
	if err != nil {
		log.Info("failed to prepare batch: %v", err)
		return err
	}
	for _, event := range eventbatch {
		err := batch.AppendStruct(&event)
		batch.Append()
		if err != nil {
			log.Info("failed to append struct: %v", err)
			return err
		}
	}
	err = batch.Send()
	if err != nil {
		log.Info("failed to send batch: %v", err)
		return err
	} // slice will be deallocated here, fix
	return nil
}

// =========================== ПОДКЛЮЧЕНИЕ НИКОГДА НЕ ЗАКРЫВАЕТСЯ, ИСПРАВИТЬ ============================
// Get conn
func Click() (driver.Conn, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
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
