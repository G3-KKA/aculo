package httpctl

import (
	"context"
	"master-service/config"
	"master-service/internal/logger"
	"master-service/internal/model"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DefaultTimeout    = time.Second * 5
	ReadHeaderTimeout = DefaultTimeout
	ShutdownTimeout   = DefaultTimeout
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname=MockService --outpkg=mock_httpctl
type Service any

// KafkaUsecase should prepare everything for future client to connect and use.
//
//go:generate mockery --filename=mock_kafka_usecase.go --name=KafkaUsecase --dir=. --structname=MockKafkaUsecase --outpkg=mock_httpctl
type KafkaUsecase interface {

	// RegisterKafkaClient prepares kafka and SA-cluster to handle new client.
	RegisterKafkaClient() (model.KafkaMetadata, error)
}

/*
ucase

m's  = metrics()
said = choseBestSA(m's)
kafkaMD = createTopic()

HandleTopic


*/

type (
	HTTPController struct {
		l      logger.Logger
		server *http.Server

		cfg config.HTTPServer
	}
	httpKafkaHandler struct {
		ucase KafkaUsecase
	}
)

// Listen to the address in config.
//
// Blocking function, will block until .Shutdown() will be called.
func (ctl *HTTPController) Serve(ctx context.Context) error {
	lis, err := net.Listen("tcp", ctl.cfg.Address)
	if err != nil {
		return err
	}

	return ctl.server.Serve(lis)
}

// Creates controller, assign routes to http server.
func NewHTTPController(cfg config.HTTPServer, l logger.Logger, srvc Service) (*HTTPController, error) {
	engine := gin.New()
	gin.DefaultWriter = &l
	khandler := &httpKafkaHandler{}
	engine.POST("register", khandler.Register)
	ctrl := &HTTPController{
		l:   l,
		cfg: cfg,
		server: &http.Server{
			Addr:                         "",
			Handler:                      engine,
			DisableGeneralOptionsHandler: false,
			TLSConfig:                    nil,
			ReadTimeout:                  ReadHeaderTimeout,
			ReadHeaderTimeout:            ReadHeaderTimeout,
			WriteTimeout:                 0,
			IdleTimeout:                  0,
			MaxHeaderBytes:               0,
			TLSNextProto:                 nil,
			ConnState:                    nil,
			ErrorLog:                     nil,
			BaseContext:                  nil,
			ConnContext:                  nil,
		},
	}

	return ctrl, nil
}

// Gracefully shutdown the server.
//
// After [ShutdownTimeout] will return [ErrShutdownTimeoutExceeded].
//
// Blocking function, will block until .Shutdown() will be called.
func (ctl *HTTPController) Shutdown(ctx context.Context) error {
	ctx = context.WithoutCancel(ctx)
	ctx, _ = context.WithTimeoutCause(ctx, ShutdownTimeout, ErrShutdownTimeoutExceeded)

	return ctl.server.Shutdown(ctx)
}

// TODO: Swagger here.
func (h *httpKafkaHandler) Register(gctx *gin.Context) {
	md , err := h.ucase.RegisterKafkaClient()
	ДАЛЬШЕ БЛЯТЬ 
	gctx.JSON(http.StatusOK, "working!")
}
