package service

import (
	"aculo/connector-restapi/internal/testutils"
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

/* func (t *testSuite) Test_Mock_SendEvent() {
	mock_repo := eventrepository.NewMockEventRepository(t.T())
	mock_repo.On("SendEvent", context.TODO(), eventrepository.SendEventRequest{
		EID: "1",
	}).Return(eventrepository.SendEventResponse{
		Event: domain.Event{
			Data: ([]byte)("[INFO] other me mario"),
		},
	}, error(nil))

	service, err := New(context.TODO(), config.Get(), BuildEserviceRequest{

		Repo: mock_repo,
	})
	t.Equal(nil, err)
	rsp, err := service.SendEvent(context.TODO(), SendEventRequest{
		EID: "1",
	})
	t.Equal(nil, err)

} */

type depInjRepository struct {
	db map[string]string
}

/* func (d *depInjRepository) SendEvent(ctx context.Context, req eventrepository.SendEventRequest) (eventrepository.SendEventResponse, error) {
	retevent := domain.Event{
		Data: []byte(d.db[req.EID]),
	}
	return eventrepository.SendEventResponse{
		Event: retevent,
	}, nil
} */

/* func (t *testSuite) Test_DependencyInjection_SendEvent() {
	DEPINJrepo := &depInjRepository{}
	DEPINJrepo.db = map[string]string{"1": "[INFO] other me mario"}
	service, err := New(context.TODO(), config.Get(), BuildEserviceRequest{

		Repo: DEPINJrepo,
	})
	t.Equal(nil, err)
	rsp, err := service.SendEvent(context.TODO(), SendEventRequest{
		EID: "1",
	})
	t.Equal(nil, err)
	panic("code below this statement broken HASH:j1kbb23jkf300")
	t.Equal(struct{}{}, rsp.Event)

}
*/
