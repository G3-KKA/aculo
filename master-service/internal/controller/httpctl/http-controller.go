package httpctl

import (
	"context"
	"master-service/internal/config"
	"master-service/internal/logger"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

// Listen to the address in config
// Synchronous function, will block untill .Shutdown() will be called
func (ctl *HTTPController) Serve(ctx context.Context) error {
	lis, err := net.Listen("tcp", ctl.cfg.Address)
	if err != nil {
		return err
	}
	return ctl.server.Serve(lis)
}

// Creates controller, assign routes to http server
func NewHTTPController(cfg config.HTTPServer, l logger.Logger, srvc Service) (*HTTPController, error) {
	engine := gin.New()
	gin.DefaultWriter = &l
	khandler := &httpKafkaHandler{}
	engine.POST("register", khandler.Register)
	ctrl := &HTTPController{
		l: l,
		// srvc: srvc,
		cfg: cfg,
		server: &http.Server{
			//			Addr:    "",
			Handler: engine,
		},
	}
	return ctrl, nil
}

const HTTPCTL_TIMEOUT = time.Second * 5

// # Synchronous function, will block
//
// Gracefully shutdown the server.
//
// After [HTTPCTL_TIMEOUT] timeout will return [ErrShutdownTimeoutExceeded]
func (ctl *HTTPController) Shutdown(ctx context.Context) error {
	ctx = context.WithoutCancel(ctx)
	ctx, _ = context.WithTimeoutCause(ctx, HTTPCTL_TIMEOUT, ErrShutdownTimeoutExceeded)
	return ctl.server.Shutdown(ctx)
}

// TODO: Swagger here
func (h *httpKafkaHandler) Register(gctx *gin.Context) {
	gctx.JSON(200, "working!")

}
