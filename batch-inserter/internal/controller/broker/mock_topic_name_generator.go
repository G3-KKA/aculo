// Code generated by mockery v2.43.2. DO NOT EDIT.

package broker

import mock "github.com/stretchr/testify/mock"

// MockTopicNameGenerator is an autogenerated mock type for the TopicNameGenerator type
type MockTopicNameGenerator struct {
	mock.Mock
}

// Generate provides a mock function with given fields:
func (_m *MockTopicNameGenerator) Generate() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Generate")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockTopicNameGenerator creates a new instance of MockTopicNameGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTopicNameGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTopicNameGenerator {
	mock := &MockTopicNameGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
