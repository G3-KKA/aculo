package httpinterface

import (
	aculo "aculo/client/pkg"
	"aculo/client/pkg/unified/unierrors"
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"

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
	logch chan []byte
	topic string

	cfg ConfigHTTP

	closer chan struct{}
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

func (h *httpapi) Close() error {
	panic("unimplemented .Close httpapi")
}

// Write implements aculo.Conn.
func (h *httpapi) Write(plainlog []byte) (n int, err error) {
	reader := bytes.NewReader(plainlog)
	url := h.cfg.core.Dst + "?topic=" + h.topic
	_, err = http.Post(url, "Content-Type:text/html; charset=UTF-8", reader)
	return
}

const noexport_ROUTINES_HANDLERS_COUNT = 1

// Zero check variant
// Returning channel never closed
func FanIn[T any](in ...<-chan T) chan T {
	ret := make(chan T, len(in))

	listener := func(in <-chan T) {

		var value T

		for {
			value = <-in
			ret <- value
		}
	}
	for i := range len(in) {
		go listener(in[i])
	}
	return ret
}

func (h *httpapi) serve(ctx context.Context) {

	errgroup := errgroup.Group{}
	errgroup.SetLimit(noexport_ROUTINES_HANDLERS_COUNT + 1)

	done := FanIn(ctx.Done(), h.closer)

	logroutine := func() error {

		var (
			err     error
			rawlog  []byte
			written int
		)

		for {
			select {
			case <-done:
				return nil
			case rawlog = <-h.logch:
			}

			written, err = h.Write(rawlog)

			if err != nil {
				return err
			}
			if written != len(rawlog) {
				return unierrors.ErrLogWrittenPartially
			}

		}
	}
	for range noexport_ROUTINES_HANDLERS_COUNT {

		errgroup.Go(logroutine)
	}
	err := errgroup.Wait()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

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

type metadata struct {
	Address string `json:"address"`
	Topic   string `json:"topic"`
}

const (
	noexport_DEFAULT_LCHANNEL_SIZE = 100
)

func NewHTTP(ctx context.Context, controllerAddr string, options ...OptionFunc) (SingleStringLogger, error) {

	// Metadata
	rsp, err := http.Get(controllerAddr + noexport_METADATA_QUERY)
	if err != nil {

		return nil, errors.Join(unierrors.ErrUnsuccessfulInitialisation, err)
	}
	defer rsp.Body.Close()

	//
	var md metadata
	body := make([]byte, rsp.ContentLength)
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
		topic: md.Topic,
		// pool:   sync.Pool{},
		cfg:    cfg,
		closer: make(chan struct{}, 1),
		logch:  make(chan []byte, noexport_DEFAULT_LCHANNEL_SIZE),
	}
	go api.serve(ctx)
	return api, nil
}

// Log implements aculo.Conn.
func (h *httpapi) Log(msg string) {
	// TODO: make use of returned
	h.logch <- []byte(msg)
}
func DefaultConfig(addr string) ConfigHTTP {
	return ConfigHTTP{
		core: aculo.ConfigCore{
			Dst: addr,
		},
	}
}
