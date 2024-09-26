package grpcctl

import (
	"context"
	"master-service/internal/config"
	"master-service/internal/logger"
	"net"

	"google.golang.org/grpc"
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname=MockService --outpkg=mock_grpcctl
type Service interface {
}
type (
	GRPCController struct {
		l   logger.Logger
		cfg config.GRPCServer

		server *grpc.Server
		kh     *grpcKafkaHandler
	}
	grpcKafkaHandler struct {
		UnimplementedRegistratorServer
	}
)

// Creates controller, registers handlers for rpc
func NewGRPCController(cfg config.GRPCServer, l logger.Logger, srvc Service) (*GRPCController, error) {
	server := grpc.NewServer()
	grpcC := &GRPCController{
		l: l,
		// srvc:             srvc,
		cfg:    cfg,
		server: server,
		kh:     &grpcKafkaHandler{},
	}
	RegisterRegistratorServer(server, grpcC.kh)
	return grpcC, nil
}

// Starts listen to the address
// Blocking execution untill .Shutdown() will be called
func (ctl *GRPCController) Serve(ctx context.Context) error {
	lis, err := net.Listen("tcp", ctl.cfg.Address)
	if err != nil {
		return err
	}
	return ctl.server.Serve(lis)
}

// Gracefully shutdown the controller
// Blocking execution untill all clients finished
func (ctl *GRPCController) Shutdown(ctx context.Context) error {
	ctl.server.GracefulStop()
	return nil
}

func (h *grpcKafkaHandler) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	panic("unimplemented grpc handler register ")
}
