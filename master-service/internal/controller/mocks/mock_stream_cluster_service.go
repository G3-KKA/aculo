// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock_controller

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	req "master-service/internal/req"
)

// MockStreamClusterService is an autogenerated mock type for the StreamClusterService type
type MockStreamClusterService struct {
	mock.Mock
}

// HandleTopic provides a mock function with given fields: ctx, _a1
func (_m *MockStreamClusterService) HandleTopic(ctx context.Context, _a1 req.MetricsRequest) (req.HandleTopicResponse, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for HandleTopic")
	}

	var r0 req.HandleTopicResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, req.MetricsRequest) (req.HandleTopicResponse, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, req.MetricsRequest) req.HandleTopicResponse); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(req.HandleTopicResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, req.MetricsRequest) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Metrics provides a mock function with given fields: ctx, _a1
func (_m *MockStreamClusterService) Metrics(ctx context.Context, _a1 req.MetricsRequest) (req.MetricsResponse, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Metrics")
	}

	var r0 req.MetricsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, req.MetricsRequest) (req.MetricsResponse, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, req.MetricsRequest) req.MetricsResponse); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(req.MetricsResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, req.MetricsRequest) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockStreamClusterService creates a new instance of MockStreamClusterService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStreamClusterService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStreamClusterService {
	mock := &MockStreamClusterService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}