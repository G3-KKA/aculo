package repository

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/generics/batchprovider"
	"aculo/batch-inserter/internal/generics/streampool"
	"aculo/batch-inserter/internal/interfaces/shuttable"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/logger"
	"context"
	"sync"
	"sync/atomic"
)

// Static check
func _() {
	var ( // [Repository]
		_ shuttable.Shuttable              = (*repository)(nil)
		_ txface.Tx[*repoapi]              = (*repository)(nil)
		_ txface.ApiWrapper[RepositoryAPI] = (*repository)(nil)
		_ Repository                       = (*repository)(nil)
	)
	var _ RepositoryAPI = (*repoapi)(nil)
}

//go:generate mockery --filename=mock_repositoryapi.go --name=RepositoryAPI --dir=. --structname MockRepositoryAPI  --inpackage=true
type RepositoryAPI interface {
	//
	// Send len(batch) to database as single query
	SendBatch(ctx context.Context, batch []domain.Log) error
	//
	// Accumulate logs til limit of its own buffer,
	// then send them to database as single query
	HandleLogStream(ctx context.Context, logs <-chan domain.Log) error
}

//go:generate mockery --filename=mock_repository.go --name=Repository --dir=. --structname MockRepository  --inpackage=true
type Repository interface {
	shuttable.Shuttable
	txface.Tx[*repoapi]
	txface.ApiWrapper[RepositoryAPI]
}

func New(ctx context.Context, config config.Config, l logger.Logger) (*repository, error) {
	conn, err := Click(config)
	if err != nil {
		return nil, err
	}

	bprov := batchprovider.New[domain.Log](ctx, config)

	api := &repoapi{
		ch:            conn,
		logger:        l,
		pool:          streampool.NewStreamPool(),
		batchprovider: bprov,
	}
	repo := &repository{
		r:           api,
		unavailable: atomic.Bool{},
		mx:          sync.RWMutex{},
	}
	return repo, nil

}
