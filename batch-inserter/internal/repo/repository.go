package repository

import "context"

//go:generate mockery --filename=mock_repository.go --name=Repository --dir=. --structname MockRepository  --inpackage=true
type Repository interface {
	SendBatch(context.Context, SendBatchRequest) error
}
type SendBatchRequest struct {
}
