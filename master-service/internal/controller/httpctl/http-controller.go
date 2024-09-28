package httpctl

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"master-service/internal/config"
	"master-service/internal/logger"
)

const (
	DEFAULT_TIMEOUT     = time.Second * 5
	READ_HEADER_TIMEOUT = DEFAULT_TIMEOUT
	SHUTDOWN_TIMEOUT    = DEFAULT_TIMEOUT
)

//go:generate mockery --filename=mock_service.go --name=Service --dir=. --structname=MockService --outpkg=mock_httpctl
type Service interface {
}
type (
	HTTPController struct {
		l      logger.Logger
		server *http.Server

		cfg config.HTTPServer
	}
	httpKafkaHandler struct {
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
			ReadTimeout:                  READ_HEADER_TIMEOUT,
			ReadHeaderTimeout:            READ_HEADER_TIMEOUT,
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
// After [SHUTDOWN_TIMEOUT] will return [ErrShutdownTimeoutExceeded].
//
// Blocking function, will block until .Shutdown() will be called.
func (ctl *HTTPController) Shutdown(ctx context.Context) error {
	ctx = context.WithoutCancel(ctx)
	ctx, _ = context.WithTimeoutCause(ctx, SHUTDOWN_TIMEOUT, ErrShutdownTimeoutExceeded)

	return ctl.server.Shutdown(ctx)
}

// TODO: Swagger here.
func (h *httpKafkaHandler) Register(gctx *gin.Context) {
	gctx.JSON(http.StatusOK, "working!")

}
