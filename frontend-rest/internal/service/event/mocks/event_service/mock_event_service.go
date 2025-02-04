// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock_event_service

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	request "aculo/frontend-restapi/internal/request"
)

// MockEventService is an autogenerated mock type for the EventService type
type MockEventService struct {
	mock.Mock
}

// GetEvent provides a mock function with given fields: ctx, req
func (_m *MockEventService) GetEvent(ctx context.Context, req request.GetEventRequest) (request.GetEventResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetEvent")
	}

	var r0 request.GetEventResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.GetEventRequest) (request.GetEventResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.GetEventRequest) request.GetEventResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(request.GetEventResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.GetEventRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockEventService creates a new instance of MockEventService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventService {
	mock := &MockEventService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
