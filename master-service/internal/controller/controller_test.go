package controller

import (
	"context"
	"testing"

	"master-service/config"
	"master-service/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	mock_controller "master-service/internal/controller/mocks"
)

func TestInitShutdown(t *testing.T) {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		egroup     errgroup.Group
		err        error
		controller *controller
		cfg        config.Controller
	)
	initialize := func() {
		cfg.GRPCServer.Address = "localhost:22222"
		cfg.HTTPServer.Address = "localhost:22223"
		ctx, cancel = context.WithCancel(context.Background())
		mockService := mock_controller.NewMockService(t)
		mockService.On("Shutdown").Return(error(nil))
		controller, err = New(cfg, logger.Noop(), mockService)

		assert.NoError(t, err)
		require.NotNil(t, controller)

		egroup.Go(func() error {
			return controller.Serve(ctx)
		})
	}
	initialize()

	cancel()
	err = egroup.Wait()
	assert.NoError(t, err)

}

/* TODO:
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
	rsp, err = http.Post(query, "application/json", &bytes.Buffer{})
	assert.NoError(t, err)

	body = make([]byte, rsp.ContentLength)
	_, err = rsp.Body.Read(body)
	assert.NoError(t, err)

	err = sonic.Unmarshal(body, &md)
	assert.NoError(t, err)

} */

/* TODO:
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
*/
