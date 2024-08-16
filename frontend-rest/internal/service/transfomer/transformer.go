package transfomer

import (
	"aculo/frontend-restapi/internal/config"
	"context"
)

type TransformRequest struct {
	SpecifiedSchema struct{}
	Data            string
}
type TransformResponse struct {
	Data struct{}
}

//go:generate mockery --name=Transformer --dir=. --outpkg=mock_transformer --filename=mock_transformer.go --output=./mocks/transformer --structname MockTransformer
type Transformer interface {
	Transform(ctx context.Context, req TransformRequest) (TransformResponse, error)
}

type transformer struct {
}

func New(ctx context.Context, config config.Config) Transformer {
	return &transformer{}
}

func (t *transformer) Transform(ctx context.Context, req TransformRequest) (TransformResponse, error) {
	return TransformResponse{}, nil
}
