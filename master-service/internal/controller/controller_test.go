package controller

import (
	"context"
	"master-service/internal/config"
	"master-service/internal/logger"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func TestInitClose(t *testing.T) {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		egroup     errgroup.Group
		err        error
		controller *controller
		config     config.Controller
	)

	ctx, cancel = context.WithCancel(context.Background())

	controller, err = New(config, logger.NewNoopLogger(), any(nil), WithContext(ctx))

	assert.NoError(t, err)

	egroup.Go(func() error {
		return controller.Serve(ctx)
	})
	cancel()
	egroup.Wait()

}

// todo
func TestRegister(t *testing.T) {

	const query = "http://localhost:7730/register"
	var (
		err error

		rsp  *http.Response
		body []byte

		md struct {
			Address string `json:"address"`
			Topic   string `json:"topic"`
		}
	)
	rsp, err = http.Get(query)
	assert.NoError(t, err)

	body = make([]byte, rsp.ContentLength)
	_, err = rsp.Body.Read(body)
	assert.NoError(t, err)

	err = sonic.Unmarshal(body, &md)
	assert.NoError(t, err)

}

// todo
func TestRegisterGRPC(t *testing.T) {
	var (
		md struct {
			Address string `json:"address"`
			Topic   string `json:"topic"`
		}
	)
	// Или как там это делается
	grpc.Get()
}
