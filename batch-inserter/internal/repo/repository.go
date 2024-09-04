package repository

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/unified/unierrors"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var _ unifaces.Tx[RepositoryAPI] = (*clickhouseRepo)(nil)
var _ Repository = (*clickhouseRepo)(nil)

//go:generate mockery --filename=mock_repository.go --name=Repository --dir=. --structname MockRepository  --inpackage=true
type Repository interface {
	unifaces.Tx[RepositoryAPI]
	RepositoryAPI
}

//go:generate mockery --filename=mock_repositoryapi.go --name=RepositoryAPI --dir=. --structname MockRepositoryAPI  --inpackage=true
type RepositoryAPI interface {
	SendBatch(ctx context.Context, batch []domain.Event) error
	GracefulShutdown() error
}
type clickhouseRepo struct {
	ch     clickhouse.Conn
	logger logger.Logger

	unavailable atomic.Bool
	mx          *sync.RWMutex
}

// GracefulShutdown implements RepositoryAPI.
func (repo *clickhouseRepo) GracefulShutdown() error {
	if repo.unavailable.CompareAndSwap(false, true) {
		// GracefulShutdown can be called only if Tx was called, so we need a bit of magic here
		repo.mx.RUnlock()
		// Lock cant be acquired while Tx with
		repo.mx.Lock()
		defer repo.mx.RLock()
		defer repo.mx.Unlock()
		return repo.ch.Close()
	}
	return unierrors.ErrUnavailable
}

// # Common middleware for all api calls
//
// Safe to call multiple times, will return [unifaces.ErrTxAlreadyClosed].
// This error may be ommited, because multi-call cannot break logic

func (repo *clickhouseRepo) Tx() (RepositoryAPI, unifaces.TxClose, error) {

	if repo.unavailable.Load() {
		return nil, func() error { return unierrors.ErrUnavailable }, unierrors.ErrUnavailable
	}

	repo.mx.RLock()
	var closed atomic.Bool
	f := func() error {
		if !closed.CompareAndSwap(false, true) {
			return unifaces.ErrTxAlreadyClosed
		}
		repo.mx.RUnlock()
		return nil
	}
	return repo, unifaces.TxClose(f), nil

}
func New(ctx context.Context, config config.Config, l logger.Logger) (*clickhouseRepo, error) {
	conn, err := Click(config)
	if err != nil {
		return nil, err
	}
	return &clickhouseRepo{
		ch:     conn,
		logger: l,
		mx:     &sync.RWMutex{},
	}, nil

}

const noexport_SEND_BATCH_QUERY = `INSERT INTO event.main_table (eid, provider_id, schema_id, type, data)`

func (c *clickhouseRepo) SendBatch(ctx context.Context, eventbatch []domain.Event) error {

	batch, err := c.ch.PrepareBatch(ctx, noexport_SEND_BATCH_QUERY)
	if err != nil {
		c.logger.Info(err)
		return err
	}
	for _, event := range eventbatch {
		err := batch.AppendStruct(&event)
		if err != nil {
			c.logger.Info(err)
			return err
		}
		/* 		err = batch.Append()
		   		if err != nil {
		   			c.logger.Info(err)
		   			return err
		   		} */
	}
	err = batch.Send()
	if err != nil {
		c.logger.Info(err)
		return err
	} // slice will be deallocated here, fix
	return nil
}

// =========================== ПОДКЛЮЧЕНИЕ НИКОГДА НЕ ЗАКРЫВАЕТСЯ, ИСПРАВИТЬ ============================
// Get conn
// КУЧА ХАРДКОДА
func Click(cfg config.Config) (driver.Conn, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: cfg.Repository.Addresses,
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
