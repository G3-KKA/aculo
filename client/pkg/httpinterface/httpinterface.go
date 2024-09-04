package httpinterface

import (
	aculo "aculo/client/pkg"
	"aculo/client/pkg/unified/unierrors"
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/errgroup"
)

var _ SingleStringLogger = &httpapi{}
var _ aculo.Conn = SingleStringLogger(nil)

//go:generate mockery --filename=mock_single_string_logger.go --name=SingleStringLogger --dir=. --structname MockSingleStringLogger  --inpackage=true
type SingleStringLogger interface {
	aculo.Conn
	Log(msg string)
}
type httpapi struct {
	topic string
	cfg   ConfigHTTP

	pool  sync.Pool
	logch chan []byte
}

/*

Что нужно сделать

Логгер ходит на *контроллера* (( ещё нет такого ))

TODO  переименовать batch-inserter ? Разбить на 2 отдельных сервиса склееных в 1 го процесс ?
FORK ? -- бессмысленно, нам хватит одного треда , пока что, простно юзаем горутины

Тот отдает метаданные (( пока что просто топик, для которого будет в последствии запущен batch-inserter ))

На этом !!! инициализация !!! New()=>httpapi заканчиваетяс

Потом в полученную метадату уже срём логами

// Горутина логгер недоделана

// Похуй на оптимизацию, для нее существует gRPC Stream и WebSocket

*/

func (h *httpapi) Close() error { return nil }

// Write implements aculo.Conn.
func (h *httpapi) Write(plainlog []byte) (n int, err error) {
	reader := bytes.NewReader(plainlog)
	_, err = http.Post(string(h.cfg.core.Dst), "Content-Type:text/html; charset=UTF-8", reader)
	return
}

const noexport_ROUTINES_HANDLERS_COUNT = 1

func (h *httpapi) serve(ctx context.Context) (err error) {
	// wg := sync.WaitGroup{}
	errgroup := errgroup.Group{}
	errgroup.SetLimit(noexport_ROUTINES_HANDLERS_COUNT + 1)
	for range noexport_ROUTINES_HANDLERS_COUNT {
		//wg.Add(1)
		f := func() error {
			for {
				select {
				case <-ctx.Done():
				case rawlog := <-h.logch:
					todo, todo2 := h.Write(rawlog)
				}

			}
		}
		errgroup.Go(f)
	}
	err2 := errgroup.Wait()
	err = errors.Join(err2, err)
	return

}

type ConfigHTTP struct {
	core aculo.ConfigCore
}
type OptionFunc func(cfg *ConfigHTTP) error

/*
	 func WithMultiDestination(dst aculo.Destination)OptionFunc {
		f := OptionFunc(func(cfg *ConfigHTTP) error {
			cfg.core.Dst

		})
		return f
	}
*/
// Заменить /metadata на controller.METADATA_QUERY
const noexport_METADATA_QUERY = "/metadata"

type noexport_MetadataResponse struct {
	address string
	topic   string
}

const noexport_DEFAULT_LCHANNEL_SIZE = 100

func NewHTTP(ctx context.Context, controllerAddr string, options ...OptionFunc) (SingleStringLogger, error) {

	// Metadata
	if !aculo.Destination(controllerAddr).ValidHTTP() {
		return nil, unierrors.ErrUnsuccessfulInitialisation
	}
	rsp, err := http.Get(controllerAddr + noexport_METADATA_QUERY)
	if err != nil {

		return nil, errors.Join(unierrors.ErrUnsuccessfulInitialisation, err)
	}

	//
	var md noexport_MetadataResponse
	body := make([]byte, rsp.ContentLength, rsp.ContentLength)
	_, err = rsp.Body.Read(body)
	if err != nil {
		return nil, errors.Join(unierrors.ErrUnsuccessfulInitialisation, err)
	}
	err = sonic.Unmarshal(body, &md)
	if err != nil {
		return nil, errors.Join(unierrors.ErrUnsuccessfulInitialisation, err)
	}

	//
	cfg := DefaultConfig(controllerAddr)
	for _, opt := range options {
		opt(&cfg)
	}

	api := &httpapi{
		topic: md.topic,
		pool:  sync.Pool{},
		cfg:   cfg,
		logch: make(chan []byte, noexport_DEFAULT_LCHANNEL_SIZE),
	}
	go api.serve(ctx)
	return api, nil
}

// Log implements aculo.Conn.
func (h *httpapi) Log(msg string) {
	// TODO: make use of returned
	_, _ = h.Write([]byte(msg))
	return
}
func DefaultConfig(addr string) ConfigHTTP {
	return ConfigHTTP{
		core: aculo.ConfigCore{
			Dst: aculo.Destination(addr),
		},
	}
}
func main() {

}
