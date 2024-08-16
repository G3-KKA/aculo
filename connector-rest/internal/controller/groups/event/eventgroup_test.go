package event

import (
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/service"
	"aculo/connector-restapi/internal/testutils"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type EventgroupTestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(EventgroupTestSuite))
}

// ========================

func (t *EventgroupTestSuite) SetupSuite() {
	testutils.DefaultSetup(t, "../../../..")
}
func (t *EventgroupTestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	}

}
func (t *EventgroupTestSuite) Test_sendSingleEvent() {
	time.Sleep(1 * time.Second)

	mock_service := service.NewMockService(t.T())
	serviceReq := service.SendEventRequest{
		Topic: "test",
		Event: []byte(`{"id": "1","name": "joe"}`),
	}
	mock_service.On("SendEvent", context.TODO(), serviceReq).Return(service.SendEventResponse{}, error(nil))

	gin.SetMode(gin.TestMode)
	testrouter := gin.Default()
	eGroup := NewEventGroup(context.TODO(), config.Config{}, mock_service)
	eGroup.Attach(&testrouter.RouterGroup)

	req, _ := http.NewRequest("POST", "/event/?topic=test", strings.NewReader(`{"id": "1","name": "joe"}`))
	w := httptest.NewRecorder()
	testrouter.ServeHTTP(w, req)
	t.Equal(w.Code, 200)
	t.Equal(w.Body.String(), `{"status":"ok"}`)

}
