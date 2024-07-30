package brocker

import (
	"aculo/batch-inserter/domain"
	"context"
)

//go:generate mockery --filename=mock_brocker.go --name=Brocker --dir=. --structname MockBrocker  --inpackage=true
type Brocker interface {
	Consume(context.Context) (domain.Event, error)
}
