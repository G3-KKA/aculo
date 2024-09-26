package controller

import (
	"context"
	"master-service/internal/config"
	mock_controller "master-service/internal/controller/mocks"
	"master-service/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func initialise() {

}
func TestInitShutdown(t *testing.T) {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		egroup     errgroup.Group
		err        error
		controller *controller
		config     config.Controller
	)
	initialise := func() {
		config.GRPCServer.Address = "localhost:22222"
		config.HTTPServer.Address = "localhost:22223"
		ctx, cancel = context.WithCancel(context.Background())
		mockService := mock_controller.NewMockService(t)
		mockService.On("Shutdown").Return(error(nil))
		controller, err = New(config, logger.Noop(), mockService)

		assert.NoError(t, err)
		require.NotNil(t, controller)

		egroup.Go(func() error {
			return controller.Serve(ctx)
		})
	}
	initialise()

	cancel()
	egroup.Wait()

}

// todo
/* func TestRegister(t *testing.T) {

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

// todo
/*
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
