package controller

import (
	"context"
	"master-service/internal/config"
	"master-service/internal/logger"
	"net/http"

	"google.golang.org/grpc"
)

type (
	controller struct {
		l    logger.ILogger
		srvc Service

		grpcServer *grpc.Server
		httpServer *http.Server
		UnimplementedRegistratorServer
	}

	grpcHandler struct {
	}
	restHandler struct {
	}
	//go:generate mockery --filename=mock_server.go --name=Service --dir=. --structname MockService  --inpackage=true
	Service interface {
		StreamClusterService
		KafkaClusterService
	}
	//go:generate mockery --filename=mock_kafka_cluster_service.go --name=KafkaClusterService --dir=. --structname MockKafkaClusterService  --inpackage=true
	KafkaClusterService interface {
	}
	//go:generate mockery --filename=mock_stream_cluster_service.go --name=StreamClusterService --dir=. --structname MockStreamClusterService  --inpackage=true
	StreamClusterService interface {
		Metrics(ctx context.Context) (*domain.Metrics, error)
		HandleTopic(ctx context.Context, said domain.SAID, topic string) error
	}
)

func New(cfg config.Controller, l logger.ILogger, srvc Service, opts ...OptionFunc) (*controller, error) {
	// lis, err := net.Listen("tcp", cfg.Address) // Это должно быть в Serve
	server := grpc.NewServer()
	ctrl := &controller{
		l:          l,
		srvc:       srvc,
		grpcServer: server,
		httpServer: &http.Server{
			Addr: cfg.HTTPServer.Address,
			/* MULIPLEXER HERE  */ Handler: nil,
		},
		UnimplementedRegistratorServer: UnimplementedRegistratorServer{},
	}
	RegisterRegistratorServer(server, ctrl)
	return ctrl, nil

}
