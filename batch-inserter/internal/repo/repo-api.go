package repository

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/generics/batchprovider"
	"aculo/batch-inserter/internal/generics/streampool"
	"aculo/batch-inserter/internal/logger"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
)

var (
	ErrReturnForced = errors.New("return while channel potentially have messages")
)

type (
	repoapi struct {
		logger logger.Logger

		ch            clickhouse.Conn
		pool          *streampool.Pool
		batchprovider *batchprovider.BatchProvider[domain.Log]
	}
)

func (repo *repoapi) shutdown() (err error) {
	repo.pool.ShutdownWait()
	err = errors.Join(err, repo.ch.Close())
	return

}

// HandleLogStream implements RepositoryAPI.
func (repo *repoapi) HandleLogStream(ctx context.Context, logs <-chan domain.Log) error {

	// repo.batchprovider.Get()
	poolFn := func(stop <-chan struct{}) {
		batch := make([]domain.Log, 10)
		var (
			idx int

			log domain.Log
			ok  bool
		)
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			case log, ok = <-logs:
				if !ok {
					return
				}
			}
			batch[idx] = log
			idx++
			if idx > 9 { // 10
				err := repo.SendBatch(ctx, batch)
				if err != nil {
					return
				}
				idx = 0
			}
		}
	}
	randomStreamHandlerName := uuid.New().String()
	repo.pool.Go(randomStreamHandlerName, poolFn)
	return nil
}

const (
	noexport_SEND_BATCH_QUERY = `INSERT INTO event.main_table (eid, provider_id, schema_id, type, data)`
)

func (repo *repoapi) SendBatch(ctx context.Context, eventbatch []domain.Log) error {

	batch, err := repo.ch.PrepareBatch(ctx, noexport_SEND_BATCH_QUERY)
	if err != nil {
		repo.logger.Info(err)
		return err
	}
	for _, event := range eventbatch {
		err := batch.AppendStruct(&event)
		if err != nil {
			repo.logger.Info(err)
			return err
		}
	}
	err = batch.Send()
	if err != nil {
		repo.logger.Info(err)
		return err
	} // slice will be deallocated here, fix
	return nil
}

// connects to clickhouse
func Click(cfg config.Repository) (driver.Conn, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: cfg.Addresses,
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
