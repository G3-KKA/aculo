package eventservice

import (
	"aculo/frontend-restapi/domain"
	"aculo/frontend-restapi/internal/config"
	repository "aculo/frontend-restapi/internal/repo"
	"aculo/frontend-restapi/internal/request"
	"aculo/frontend-restapi/internal/service/transfomer"
	testutils "aculo/frontend-restapi/internal/tests"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(testSuite))
}
func (t *testSuite) SetupSuite() {
	testutils.DefaultPreTestSetup("../../..")
}

func (t *testSuite) Test_Mock_GetEvent() {
	mock_repo := repository.NewMockRepository(t.T())
	mock_repo.On("GetEvent", context.TODO(), request.GetEventRequest{
		EID: "1",
	}).Return(request.GetEventResponse{
		Event: domain.Event{
			Data: ("[INFO] other me mario"),
		},
	}, error(nil))
	mock_transformer := transfomer.NewMockTransformer(t.T())
	mock_transformer.On("Transform", context.TODO(), transfomer.TransformRequest{
		SpecifiedSchema: struct{}{},
		Data:            ("[INFO] other me mario"),
	}).Return(transfomer.TransformResponse{
		Data: struct{}{},
	}, error(nil))

	service, err := New(context.TODO(), config.Get(), BuildEserviceRequest{

		Repo:        mock_repo,
		Transformer: mock_transformer,
	})
	t.NoError(err)
	_, err = service.GetEvent(context.TODO(), request.GetEventRequest{
		EID: "1",
	})

	t.NoError(err)
	/* t.Equal(struct{}{}, rsp.Event.(transfomer.TransformResponse).Data) */

}

type depInjRepository struct {
	db map[string]string
}

func (d *depInjRepository) GetEvent(ctx context.Context, req request.GetEventRequest) (request.GetEventResponse, error) {
	retevent := domain.Event{
		Data: (d.db[req.EID]),
	}
	return request.GetEventResponse{
		Event: retevent,
	}, nil
}

func (t *testSuite) Test_DependencyInjection_GetEvent() {
	DEPINJrepo := &depInjRepository{}
	DEPINJrepo.db = map[string]string{"1": "[INFO] other me mario"}
	mock_transformer := transfomer.NewMockTransformer(t.T())
	mock_transformer.On("Transform", context.TODO(), transfomer.TransformRequest{
		SpecifiedSchema: struct{}{},
		Data:            ("[INFO] other me mario"),
	}).Return(transfomer.TransformResponse{
		Data: struct{}{},
	}, error(nil))
	service, err := New(context.TODO(), config.Get(), BuildEserviceRequest{

		Repo:        DEPINJrepo,
		Transformer: mock_transformer,
	})
	t.NoError(err)
	_, err = service.GetEvent(context.TODO(), request.GetEventRequest{
		EID: "1",
	})
	t.NoError(err)

	/* t.Equal(struct{}{}, rsp.Event.(transfomer.TransformResponse).Data) */

}
