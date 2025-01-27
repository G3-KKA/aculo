package controller

import (
	"context"
	"errors"
	"master-service/config"
	"master-service/internal/controller/grpcctl"
	"master-service/internal/controller/httpctl"
	"master-service/internal/errspec"
	"master-service/internal/logger"
	"master-service/internal/req"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// controller will shutdown on any of this signals.
var SHUTDOWN_SIGNALS = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

type controller struct {
	grpcC *grpcctl.GRPCController
	httpC *httpctl.HTTPController
	srvc  Service

	l logger.Logger

	shutdownOnce sync.Once
}

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname=MockService --outpkg=mock_controller
type Service interface {
	StreamClusterService
	KafkaClusterService
	Shutdown() error
}

//go:generate mockery --filename=mock_kafka_cluster.go --name=KafkaClusterService --dir=. --structname=MockKafkaClusterService --outpkg=mock_controller
type KafkaClusterService interface {
	CreateTopic(ctx context.Context, r req.CreateTopicRequest) (req.CreateTopicResponse, error)
}

//go:generate mockery --filename=mock_stream_cluster_service.go --name=StreamClusterService --dir=. --structname=MockStreamClusterService --outpkg=mock_controller
type StreamClusterService interface {
	//  Deprecated: request response semantic will be better.
	//	Metrics(ctx context.Context) (*domain.StreamMetrics, error);
	//	HandleTopic(ctx context.Context, said domain.SAID, topic string) error;

	Metrics(ctx context.Context, r req.MetricsRequest) (req.MetricsResponse, error)

	HandleTopic(ctx context.Context, r req.MetricsRequest) (req.HandleTopicResponse, error)
}

// Creates controller, validates config, not starting to serve.
func New(cfg config.Controller, l logger.Logger, srvc Service) (ctl *controller, err error) {
	err = valid(cfg)
	if err != nil {
		return nil, err
	}

	grpcC, err := grpcctl.NewGRPCController(cfg.GRPCServer, l, srvc)
	if err != nil {
		return nil, err
	}
	httpC, err := httpctl.NewHTTPController(cfg.HTTPServer, l, srvc)
	if err != nil {
		return nil, err
	}

	ctrl := &controller{
		grpcC:        grpcC,
		httpC:        httpC,
		srvc:         srvc,
		l:            l,
		shutdownOnce: sync.Once{},
	}
	l.Debug("controller construction succeeded")

	return ctrl, nil
}

// # Serve starts grpc and http controller.
//
// Shutdown via ctx.Done or on any of [SHUTDOWN_SIGNALS].
func (ctl *controller) Serve(ctx context.Context) error {
	var (
		egroup errgroup.Group
	)
	shutdowner := func() (err error) {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, SHUTDOWN_SIGNALS...)
		select {
		case <-sigchan:
		case <-ctx.Done():
		}

		return ctl.Shutdown(ctx)
	}
	egroup.Go(func() error {
		return ctl.grpcC.Serve(ctx)
	})
	egroup.Go(func() error {
		return ctl.httpC.Serve(ctx)
	})
	egroup.Go(shutdowner)

	return egroup.Wait()
}

// Shutdown is graceful and cooperative.
//
// Blocking function.
//
// Safe to call it multiple places.
func (ctl *controller) Shutdown(ctx context.Context) (err error) {
	ctl.shutdownOnce.Do(func() {
		const allroutines = 3
		errs := make(chan error, allroutines)
		wg := sync.WaitGroup{}

		// Кооперативное выключение!
		// Обязательно ждём чтобы все закончили работу.
		// Можем встать на дедлоке здесь,
		//	 если одна из функций дурит внутри .Shutdown() контроллеров.
		// time.Timer() может быть необходим с чем-то наподобие 30 сек.

		grpcShutdowner := func() {
			defer wg.Done()
			errs <- ctl.grpcC.Shutdown(ctx)
		}
		httpShutdowner := func() {
			defer wg.Done()
			errs <- ctl.httpC.Shutdown(ctx)
		}
		errscloser := func() {
			defer close(errs)
			wg.Wait()
			errs <- ctl.srvc.Shutdown()
		} Вот это должно быть в App.Shutdown(), убрать ! 

		wg.Add(1)
		go grpcShutdowner()
		wg.Add(1)
		go httpShutdowner()

		go errscloser()
		for e := range errs {
			err = errors.Join(e, err)
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func valid(cfg config.Controller) (err error) {
	if cfg.GRPCServer.Address == cfg.HTTPServer.Address {
		err = errspec.Same(ErrConfigSameAddresses,
			cfg.GRPCServer.Address,
			cfg.HTTPServer.Address)
	}

	return
}
