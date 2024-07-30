package repository

import (
	"aculo/frontend-restapi/domain"
	"aculo/frontend-restapi/internal/config"
	"context"

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
	GetEvent(context.Context, GetEventRequest) (GetEventResponse, error)
}

func New(ctx context.Context, config config.Config, ch clickhouse.Conn) (Repository, error) {
	repo := &eRepo{conn: ch}
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
